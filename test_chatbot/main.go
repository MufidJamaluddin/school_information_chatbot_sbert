package main

import (
	abdto "chatbot_be_go/src/application/abbreviation/dto"
	dm "chatbot_be_go/src/domain"
	"chatbot_be_go/src/persistence"
	appConf "chatbot_be_go/src/persistence/config"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

	nFile, err := os.Open("data.csv")
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

		if _, err = persistenceObj.QuestionRepository.SaveNewQuestion(
			ctx,
			question,
			answer,
			roleGroupId,
			dataOwnerUserName,
		); err != nil {
			log.Printf("\nError in Save New Question for No %s: %+v\n", no, err)
		}
	}
}

// Melakukan Testing atas Data Pertanyaan dan Jawaban
func loadTestQuestionAndAnswer(
	scenarioNo uint64,
	scenarioName string,
	persistenceObj *persistence.Persistence,
) {
	// File Test Reader
	nTestDataFile, err := os.Open("test.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer nTestDataFile.Close()

	// File Test Result
	nTestResultFile, err := os.Create(
		fmt.Sprintf("test_sbert_%d_%s.csv", scenarioNo, scenarioName),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer nTestResultFile.Close()

	// Start Testing
	ctx := context.Background()

	testDataReader := csv.NewReader(nTestDataFile)
	testDataReader.Comma = ';'

	testResultWriter := csv.NewWriter(nTestResultFile)
	testResultWriter.Comma = ';'
	for {
		testDataRecord, err := testDataReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		no := testDataRecord[0]
		question := testDataRecord[1]

		answer, similarityValue, err := persistenceObj.QuestionRepository.FindAnswerWithSimilarityValue(
			ctx,
			question,
		)
		if err == sql.ErrNoRows || similarityValue < 0 {
			continue
		}
		if err != nil {
			log.Printf("\nError in Get Answer for No %s: %+v\n", no, err)
		}

		for i := len(testDataRecord); i < 4; i++ {
			testDataRecord = append(testDataRecord, "")
		}

		testDataRecord[2] = answer
		testDataRecord[3] = fmt.Sprintf("%f", similarityValue)

		testResultWriter.Write(testDataRecord)
		testResultWriter.Flush()
	}

	if err = testResultWriter.Error(); err != nil {
		log.Println("Error in Save Test Result: ", err)
	}
}

// Menjalankan Semua Proses
func main() {
	_ = godotenv.Load()
	logger := logrus.New()
	config := appConf.New()

	testScenario := []testModel{
		{
			ScenarioNo:    1,
			ScenarioName:  "ori_model",
			ModelPath:     "/ori_model",
			ModelFileName: "model.pkl",
		},
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

		log.Println("Uji Pertanyaan dan Jawaban")
		loadTestQuestionAndAnswer(scenario.ScenarioNo, scenario.ScenarioName, persistenceObj)
	}
}
