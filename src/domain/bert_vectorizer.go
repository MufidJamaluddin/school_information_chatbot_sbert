package domain

import (
	"fmt"
	"log"

	"github.com/MufidJamaluddin/transformer"
	"github.com/MufidJamaluddin/transformer/bert"
	"github.com/MufidJamaluddin/transformer/util"
	"github.com/sugarme/gotch"
	"github.com/sugarme/gotch/nn"
	"github.com/sugarme/gotch/ts"
	"github.com/sugarme/tokenizer"
)

type ISBertVectorizer interface {
	Encode(text string) ([]float64, error)
}

type sbertVectorizer struct {
	compileMode gotch.Device
	config      bert.BertConfig
	model       bert.BertEmbedding
	tokenizer   *bert.Tokenizer
}

func NewSBertVectorizer(filePath, fileName string) ISBertVectorizer {
	var config bert.BertConfig

	device := gotch.CPU
	varStore := nn.NewVarStore(device)

	_, err := util.CachedPath(filePath, fileName)
	if err != nil {
		log.Fatal("Error on loading BERT Model File", err)
	}

	if err := transformer.LoadConfig(&config, filePath, nil); err != nil {
		log.Fatal("Error on loading BERT Config Model", err)
	}

	model := bert.NewBertEmbeddings(varStore.Root(), &config)

	var tk *bert.Tokenizer = bert.NewTokenizer()
	if err := tk.Load(filePath, nil); err != nil {
		log.Fatal("Error on loading BERT Tokenizer Model", err)
	}

	return &sbertVectorizer{
		compileMode: device,
		config:      config,
		model:       model,
		tokenizer:   tk,
	}
}

func (s *sbertVectorizer) Encode(text string) (result []float64, err error) {
	singleEncodeInput := tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(text))

	encodingResult, err := s.tokenizer.Encode(singleEncodeInput, true)
	if err != nil {
		return
	}

	var attentionMask []float64
	var tokenIds []int
	var typeIds []int
	var positionIds []int

	for i := 0; i < len(encodingResult.Ids); i++ {
		if encodingResult.Words[i] < 0 {
			positionIds = append(positionIds, 0)
		} else {
			positionIds = append(positionIds, encodingResult.Words[i])
		}

		attentionMask = append(attentionMask, float64(encodingResult.AttentionMask[i]))
		tokenIds = append(tokenIds, encodingResult.Ids[i])
		typeIds = append(typeIds, encodingResult.TypeIds[i])
	}

	tokenIdsTensor := ts.TensorFrom(tokenIds).MustTo(s.compileMode, true)
	typeIdsTensor := ts.TensorFrom(typeIds).MustTo(s.compileMode, true)
	positionIdsTensor := ts.TensorFrom(positionIds).MustTo(s.compileMode, true)

	ts.NoGrad(func() {
		sbertResult, err := s.model.ForwardT(
			tokenIdsTensor,
			typeIdsTensor,
			positionIdsTensor,
			ts.None,
			false,
		)
		if err != nil {
			return
		}

		meanPoolResult, err := s.MeanPooling(sbertResult, ts.TensorFrom(attentionMask))
		if err != nil {
			return
		}

		result = meanPoolResult.Float64Values(true)
	})

	return
}

func (s *sbertVectorizer) MeanPooling(
	modelOutput *ts.Tensor,
	attentionMask *ts.Tensor,
) (*ts.Tensor, error) {
	squeezeMask, err := attentionMask.Unsqueeze(-1, true)
	if err != nil {
		fmt.Println("Error in Squeeze Mask Expanded", err)
		return nil, err
	}

	expandedMask, err := squeezeMask.Expand(modelOutput.MustSize(), false, true)
	if err != nil {
		fmt.Println("Error in Expanded Mask Expanded", err)
		return nil, err
	}

	multiple, err := modelOutput.Multiply(expandedMask, false)
	if err != nil {
		fmt.Println("Error in Multiple Embedding Value & Mask Expanded", err)
		return nil, err
	}

	multipleSum, err := multiple.SumToSize([]int64{s.config.HiddenSize}, true)
	if err != nil {
		fmt.Println("Error in Sum Multiple Embedding Value & Mask Expanded", err)
		return nil, err
	}

	expandedMaskSum, err := expandedMask.Sum(gotch.BFloat16, true)
	if err != nil {
		fmt.Println("Error in Sum Expanded Mask Expanded", err)
		return nil, err
	}

	expandedMaskClamp, err := expandedMaskSum.ClampMin(ts.FloatScalar(1e-9), true)
	if err != nil {
		fmt.Println("Error in Clamp Expanded Mask", err)
		return nil, err
	}

	return multipleSum.Divide(expandedMaskClamp, true)
}
