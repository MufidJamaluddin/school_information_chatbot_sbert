package main

import (
	abdto "chatbot_be_go/src/application/abbreviation/dto"
	dm "chatbot_be_go/src/domain"
	"chatbot_be_go/src/persistence"
	appConf "chatbot_be_go/src/persistence/config"
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type testModel struct {
	ScenarioNo    uint64
	ScenarioName  string
	ModelPath     string
	ModelFileName string
}

// Menyimpan Data Singkatan dari CSV kedalam Database
func loadAbbreviation(
	persistenceObj *persistence.Persistence,
) {
	dataOwnerUserName := "nzank"

	nAbbreviationFile, err := os.Open("abbreviation.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer nAbbreviationFile.Close()

	ctx := context.Background()

	isFirst := true

	csvReader := csv.NewReader(nAbbreviationFile)
	csvReader.Comma = ';'
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if len(rec) < 2 {
			continue
		}

		if isFirst {
			isFirst = false
			continue
		}

		formal := rec[0]
		abbreviations := strings.Split(strings.ReplaceAll(rec[1], ";", ","), ",")

		abbreviationData, errGet := persistenceObj.AbbreviationRepository.GetAbbreviation(
			ctx,
			formal,
		)

		abbreviationData.ListAbbreviationTerm = append(abbreviationData.ListAbbreviationTerm, abbreviations...)

		if errGet == nil || len(abbreviationData.ListAbbreviationTerm) > 1 {
			updatedData := &abdto.UpdateAbbreviationDTO{
				StandardWord:         abbreviationData.StandardWord,
				ListAbbreviationTerm: abbreviationData.ListAbbreviationTerm,
				UpdatedBy:            dataOwnerUserName,
			}

			if err = persistenceObj.AbbreviationRepository.UpdateAbbreviation(
				ctx,
				updatedData,
			); err != nil {
				log.Printf("\nError in Update Abbreviation Data %s: %+v\n", updatedData.StandardWord, err)
			}
		} else {
			createdData := &abdto.CreateAbbreviationDTO{
				StandardWord:         formal,
				ListAbbreviationTerm: abbreviationData.ListAbbreviationTerm,
				CreatedBy:            dataOwnerUserName,
			}

			if err = persistenceObj.AbbreviationRepository.SaveNewAbbreviation(
				ctx,
				createdData,
			); err != nil {
				log.Printf("\nError in Create Abbreviation Data %s: %+v\n", createdData.StandardWord, err)
			}
		}
	}
}

// Menyimpan Data Pertanyaan dan Jawaban dari CSV kedalam Database
func loadDataQuestionAndAnswer(
	persistenceObj *persistence.Persistence,
) {
	roleGroupId := uint64(1)
	dataOwnerUserName := "nzank"

	var wg sync.WaitGroup

	nFile, err := os.Open("load_test_data.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer nFile.Close()

	ctx := context.Background()

	if err = persistenceObj.QuestionRepository.TruncateQuestion(ctx); err != nil {
		log.Fatal(err)
		return
	}

	csvReader := csv.NewReader(nFile)
	csvReader.Comma = ';'
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if len(rec) < 3 {
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		no := rec[0]
		question := rec[1]
		answer := rec[2]

		go func() {
			wg.Add(1)
			defer wg.Done()

			if _, mErr := persistenceObj.QuestionRepository.SaveNewQuestionWithoutSBERTVector(
				ctx,
				question,
				answer,
				roleGroupId,
				dataOwnerUserName,
			); mErr != nil {
				log.Printf("\nError in Save New Question for No %s: %+v\n", no, err)
			}
		}()
	}

	wg.Wait()

	if err = persistenceObj.QuestionRepository.ResetSBERTVectorQuestion(ctx); err != nil {
		log.Printf("\nError in Reset SBERT Vector for No %s: %+v\n", no, err)
	}
}

// Menjalankan Semua Proses
func main() {
	_ = godotenv.Load()
	logger := logrus.New()
	config := appConf.New()

	testScenario := []testModel{
		{
			ScenarioNo:    2,
			ScenarioName:  "fine_tuned_model",
			ModelPath:     "/fine_tuned_model",
			ModelFileName: "model.pkl",
		},
	}

	for _, scenario := range testScenario {
		vectorizer := dm.NewSBertVectorizer(scenario.ModelPath, scenario.ModelFileName)

		persistenceObj := persistence.New(
			vectorizer,
			config.SqlDb,
			logger,
		)

		log.Println("Baca Singkatan")
		loadAbbreviation(persistenceObj)

		log.Println("Baca Semua Data Latih Pertanyaan dan Jawaban")
		loadDataQuestionAndAnswer(persistenceObj)
	}
}
