package domain

import (
	"fmt"
	"log"
	"math"

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
	angleRates  []float64
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

	angleRates := computeAngleRates(int(config.HiddenSize))

	return &sbertVectorizer{
		compileMode: device,
		config:      config,
		model:       model,
		tokenizer:   tk,
		angleRates:  angleRates,
	}
}

func (s *sbertVectorizer) Encode(text string) (result []float64, err error) {
	singleEncodeInput := tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(text))
	scale := 20

	encodingResult, err := s.tokenizer.Encode(singleEncodeInput, true)
	if err != nil {
		return
	}

	var attentionMask []float64
	for _, attentionMaskItem := range encodingResult.AttentionMask {
		attentionMask = append(attentionMask, float64(attentionMaskItem))
	}

	tokenIdsTensor := ts.TensorFrom(encodingResult.Ids).MustTo(s.compileMode, true)
	typeIdsTensor := ts.TensorFrom(encodingResult.TypeIds).MustTo(s.compileMode, true)
	positionalEncoding := createPositionalEncodingTensor(encodingResult.Words, s.angleRates, scale).MustTo(s.compileMode, true)

	ts.NoGrad(func() {
		sbertResult, err := s.model.ForwardT(
			tokenIdsTensor,
			typeIdsTensor,      // [m] Apply With [512, 768]
			positionalEncoding, // [m, 768] Apply With [2, 768]
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

	expandedMask, err := squeezeMask.Expand([]int64{modelOutput.MustSize()[0], 1}, false, true)
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

	dividedRes, err := multipleSum.Divide(expandedMaskClamp, true)
	if err != nil {
		fmt.Println("Error in Divide Multiple by Sum Mask", err)
		return nil, err
	}

	return dividedRes.Clamp(ts.FloatScalar(-1), ts.FloatScalar(1), true)
}

func createPositionalEncodingTensor(wordIndexs []int, angleRates []float64, scale int) *ts.Tensor {
	var positionalEncoding []int
	mScale := float64(scale)

	for i := 0; i < len(wordIndexs); i++ {
		positionalEncoding = append(
			positionalEncoding,
			int(math.Round(calculatePositionalEncodingItem(i, wordIndexs[i], angleRates)*mScale)),
		)
	}

	return ts.TensorFrom(positionalEncoding)
}

func calculatePositionalEncodingItem(index int, modelIndex int, angleRates []float64) float64 {
	scale := 10000.0
	angle := float64(modelIndex) / scale

	if index%2 == 0 {
		return math.Sin(angle * angleRates[index/2])
	}

	return math.Cos(angle * angleRates[index/2])
}

func computeAngleRates(embeddingSize int) []float64 {
	scale := 10000.0
	angleRates := make([]float64, embeddingSize/2)

	for i := 0; i < embeddingSize/2; i++ {
		angleRates[i] = math.Pow(scale, float64(2*i)/float64(embeddingSize))
	}

	return angleRates
}
