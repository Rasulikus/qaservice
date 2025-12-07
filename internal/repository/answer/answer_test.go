package answer

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Rasulikus/qaservice/internal/model"
	"github.com/Rasulikus/qaservice/internal/repository/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var testDB *testdb.TestDB

func TestMain(m *testing.M) {
	testDB = testdb.NewTestDB(nil, "")
	testDB.RecreateTables()
	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

type testSuite struct {
	db         *gorm.DB
	answerRepo *Repository
	ctx        context.Context
}

func setupTestSuite(t *testing.T) *testSuite {
	t.Helper()
	suite := new(testSuite)
	suite.db = testDB.DB
	suite.ctx = context.Background()
	suite.answerRepo = NewRepository(suite.db)
	testDB.CleanDB()
	return suite
}

func createQuestion(t *testing.T, db *gorm.DB) *model.Question {
	t.Helper()
	q := &model.Question{
		Text: "test question",
	}
	require.NoError(t, db.Create(q).Error)
	return q
}

func Test_Repo_Insert(t *testing.T) {
	ts := setupTestSuite(t)

	q := createQuestion(t, ts.db)
	a := &model.Answer{
		QuestionID: q.ID,
		UserID:     "user-1",
		Text:       "test answer",
	}

	testCases := []struct {
		name    string
		a       *model.Answer
		wantErr bool
	}{
		{"insert answer", a, false},
		{"insert empty answer", a, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ts.answerRepo.Create(ts.ctx, a)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				assert.NotZero(t, a.ID)
				assert.Equal(t, q.ID, a.QuestionID)
				assert.WithinDuration(t, time.Now(), a.CreatedAt, time.Second)
				require.NoError(t, err)
			}
		})
	}
}

func Test_Repo_GetByID(t *testing.T) {
	ts := setupTestSuite(t)

	q := createQuestion(t, ts.db)
	a := &model.Answer{
		QuestionID: q.ID,
		UserID:     "user-1",
		Text:       "test answer",
	}
	err := ts.answerRepo.Create(ts.ctx, a)
	require.NoError(t, err)

	t.Run("get answer by id", func(t *testing.T) {
		get, err := ts.answerRepo.GetByID(ts.ctx, a.ID)
		require.NoError(t, err)
		assert.Equal(t, a.ID, get.ID)
		assert.Equal(t, a.QuestionID, get.QuestionID)
		assert.Equal(t, a.Text, get.Text)
		assert.Equal(t, a.UserID, get.UserID)
	})

	t.Run("get answer by id not found", func(t *testing.T) {
		_, err := ts.answerRepo.GetByID(ts.ctx, -1)
		require.ErrorIs(t, err, model.ErrNotFound)
	})
}

func Test_Repo_Delete(t *testing.T) {
	ts := setupTestSuite(t)

	q := createQuestion(t, ts.db)
	a := &model.Answer{
		QuestionID: q.ID,
		UserID:     "user-1",
		Text:       "test answer",
	}
	err := ts.answerRepo.Create(ts.ctx, a)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"delete answer", a.ID, false},
		{"delete a non-existent answer", -1, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ts.answerRepo.Delete(ts.ctx, tc.id)
			if tc.wantErr {
				require.ErrorIs(t, err, model.ErrNotFound)
			} else {
				require.NoError(t, err)
				get, err := ts.answerRepo.GetByID(ts.ctx, tc.id)
				require.Error(t, err)
				require.Nil(t, get)
			}
		})
	}
}
