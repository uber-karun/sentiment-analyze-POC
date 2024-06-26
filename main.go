package main

import (
	"context"
	"fmt"
	"log"

	language "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/language/apiv1/languagepb"
	"google.golang.org/api/option"
)

func main() {
	// Load the credentials from your service account key JSON file
	ctx := context.Background()
	client, err := language.NewClient(ctx, option.WithCredentialsFile("cred.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	text := "Help, I need to contact the police. I am being followed by a strange and I am scared"

	sentiment, err := analyzeSentiment(ctx, client, text)
	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}

	fmt.Printf("Sentiment magnitude: %.2f\n", sentiment.Magnitude)

	annotation, err := annotateText(ctx, client, text)
	if err != nil {
		log.Fatalf("Failed to annotate text: %v", err)
	}

	// Display the results
	fmt.Println("result", len(annotation.Categories))
	for _, a := range annotation.Categories {
		fmt.Println("confidence", a.Name, a.Confidence)
	}
}

func analyzeSentiment(ctx context.Context, client *language.Client, text string) (*languagepb.Sentiment, error) {
	req := &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type:     languagepb.Document_PLAIN_TEXT,
			Language: "en",
		},
		EncodingType: languagepb.EncodingType_UTF8,
	}

	resp, err := client.AnalyzeSentiment(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, a := range resp.Sentences {
		fmt.Println("text:", a.Text.Content, a.Text.BeginOffset, a.Sentiment.Magnitude, a.Sentiment.Score)
	}
	return resp.DocumentSentiment, nil
}

func annotateText(ctx context.Context, client *language.Client, text string) (*languagepb.AnnotateTextResponse, error) {
	req := &languagepb.AnnotateTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		Features: &languagepb.AnnotateTextRequest_Features{
			ExtractSyntax:            true,
			ExtractEntities:          true,
			ExtractDocumentSentiment: true,
			ExtractEntitySentiment:   true,
			ClassifyText:             true,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	}
	resp, err := client.AnnotateText(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
