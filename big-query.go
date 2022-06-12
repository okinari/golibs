package golibs

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type BigQuery struct {
	projectID string
	datasetID string
	client    *bigquery.Client
	dataset   *bigquery.Dataset
}

type BigQueryTable interface {
	Error() string
}

func NewBigQueryApi(ctx context.Context, ProjectID string, DatasetID string) (*BigQuery, error) {
	bq := new(BigQuery)
	bq.projectID = ProjectID
	bq.datasetID = DatasetID

	// ctx := context.Background()
	client, err := bigquery.NewClient(ctx, bq.projectID)
	if err != nil {
		return nil, fmt.Errorf("new client create error: %v", err)
	}

	bq.client = client
	bq.dataset = client.Dataset(bq.datasetID)

	return bq, nil
}

func (bq *BigQuery) Query(ctx context.Context, sql string) [][]bigquery.Value {
	q := bq.client.Query(sql)
	it, err := q.Read(ctx)
	if err != nil {
		fmt.Print("read error.")
		return nil
	}

	var retval [][]bigquery.Value
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			fmt.Print("break.")
			break
		}
		if err != nil {
			// TODO: Handle error.
			fmt.Print("next nil")
		}
		retval = append(retval, values)
		fmt.Println(values)
	}
	return retval
}

type OptionsInsertFromCsv func(*privateOptionsInsertFromCsv)

type privateOptionsInsertFromCsv struct {
	schema             bigquery.Schema
	numSkipLeadingRows int64
}

func SetSchemaOptionsInsertFromCsv(schema bigquery.Schema) OptionsInsertFromCsv {
	return func(option *privateOptionsInsertFromCsv) {
		option.schema = schema
	}
}

func SetNumSkipLeadingRowsOptionsInsertFromCsv(num int64) OptionsInsertFromCsv {
	return func(option *privateOptionsInsertFromCsv) {
		option.numSkipLeadingRows = num
	}
}

func (bq *BigQuery) InsertFromCsv(ctx context.Context, tableName string, gcsBucketName string, gcsObjectPath string, optionsInsertFromCsv ...OptionsInsertFromCsv /*, schema bigquery.Schema*/) error {

	// 引数取得
	privateOptions := privateOptionsInsertFromCsv{}
	for _, option := range optionsInsertFromCsv {
		option(&privateOptions)
	}

	gcsRef := bigquery.NewGCSReference("gs://" + gcsBucketName + "/" + gcsObjectPath)
	gcsRef.SourceFormat = bigquery.CSV
	// gcsRef.AllowJaggedRows = true
	if privateOptions.schema != nil {
		gcsRef.AutoDetect = false
		gcsRef.Schema = privateOptions.schema
	} else {
		gcsRef.AutoDetect = true
	}

	gcsRef.SkipLeadingRows = privateOptions.numSkipLeadingRows

	loader := bq.dataset.Table(tableName).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteAppend
	// job, err = loader.Run(ctx)
	// Poll the job for completion if desired, as above.

	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}

	if status.Err() != nil {
		return fmt.Errorf("job completed with error: %v", status.Err())
	}
	return nil
}
