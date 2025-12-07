package question

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
	db           *gorm.DB
	questionRepo *Repository
	ctx          context.Context
}

func setupTestSuite(t *testing.T) *testSuite {
	t.Helper()
	suite := new(testSuite)
	suite.db = testDB.DB
	suite.ctx = context.Background()
	suite.questionRepo = NewRepository(suite.db)
	testDB.CleanDB()
	return suite
}

func createAnswer(t *testing.T, db *gorm.DB, questionID int) *model.Answer {
	t.Helper()
	a := &model.Answer{
		QuestionID: questionID,
		UserID:     "user-1",
		Text:       "test answer",
	}
	err := db.Create(a).Error
	require.NoError(t, err)

	return a
}

func Test_Repo_Create(t *testing.T) {
	ts := setupTestSuite(t)

	q := &model.Question{
		Text: "test question",
	}
	t.Run("create question", func(t *testing.T) {
		err := ts.questionRepo.Create(ts.ctx, q)
		require.NoError(t, err)
		assert.NotZero(t, q.ID)
		assert.WithinDuration(t, time.Now(), q.CreatedAt, time.Second)
	})
}

func Test_Repo_GetByID(t *testing.T) {
	ts := setupTestSuite(t)

	q := &model.Question{
		Text: "test question",
	}
	err := ts.questionRepo.Create(ts.ctx, q)
	require.NoError(t, err)

	a1 := createAnswer(t, ts.db, q.ID)
	a2 := createAnswer(t, ts.db, q.ID)

	t.Run("get question by id", func(t *testing.T) {
		get, err := ts.questionRepo.GetByID(ts.ctx, q.ID)
		require.NoError(t, err)
		assert.Equal(t, q.ID, get.ID)
		assert.Equal(t, q.Text, get.Text)

		assert.Len(t, get.Answers, 2)
		assert.Equal(t, a1.ID, get.Answers[0].ID)
		assert.Equal(t, a2.ID, get.Answers[1].ID)
	})

	t.Run("get question by id not found", func(t *testing.T) {
		_, err := ts.questionRepo.GetByID(ts.ctx, -1)
		require.ErrorIs(t, err, model.ErrNotFound)
	})
}

func Test_Repo_List(t *testing.T) {
	ts := setupTestSuite(t)

	err := ts.questionRepo.Create(ts.ctx, &model.Question{
		Text: "test question 1",
	})
	require.NoError(t, err)
	err = ts.questionRepo.Create(ts.ctx, &model.Question{
		Text: "test question 2",
	})
	require.NoError(t, err)

	t.Run("list questions", func(t *testing.T) {
		qs, err := ts.questionRepo.List(ts.ctx)
		require.NoError(t, err)
		require.Len(t, qs, 2)
	})

}

func Test_Repo_Delete(t *testing.T) {
	ts := setupTestSuite(t)

	q := &model.Question{
		Text: "test question",
	}
	err := ts.questionRepo.Create(ts.ctx, q)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"delete question", q.ID, false},
		{"delete a non-existent question", -1, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ts.questionRepo.Delete(ts.ctx, tc.id)
			if tc.wantErr {
				require.ErrorIs(t, err, model.ErrNotFound)
			} else {
				require.NoError(t, err)
				get, err := ts.questionRepo.GetByID(ts.ctx, tc.id)
				require.Error(t, err)
				require.Nil(t, get)
			}
		})
	}
}
