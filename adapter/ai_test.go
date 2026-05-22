package adapter

import "testing"

func TestNewAIAssistantOperation(t *testing.T) {
	op := NewAIAssistantOperation()
	if op == nil {
		t.Fatal("expected operation")
	}
	if op.Description == "" {
		t.Fatal("expected operation description")
	}
	if op.AdditionalProperties["capability"] != AIAssistantCapability {
		t.Fatalf("expected capability %q, got %q", AIAssistantCapability, op.AdditionalProperties["capability"])
	}
	if op.AdditionalProperties["mode"] != "read-only" {
		t.Fatalf("expected read-only mode, got %q", op.AdditionalProperties["mode"])
	}
}

func TestAIAssistantRequestValidate(t *testing.T) {
	if err := (AIAssistantRequest{}).Validate(); err == nil {
		t.Fatal("expected missing user intent error")
	}
	if err := (AIAssistantRequest{UserIntent: "explain this design"}).Validate(); err != nil {
		t.Fatalf("expected valid request, got %v", err)
	}
}

func TestAIAssistantResponseValidate(t *testing.T) {
	tests := []struct {
		name    string
		resp    AIAssistantResponse
		wantErr bool
	}{
		{
			name:    "empty response",
			resp:    AIAssistantResponse{},
			wantErr: true,
		},
		{
			name:    "explanation",
			resp:    AIAssistantResponse{Explanation: "This design contains one service."},
			wantErr: false,
		},
		{
			name: "recommendation",
			resp: AIAssistantResponse{Recommendations: []AIAssistantRecommendation{{
				Title: "Review service selectors",
			}}},
			wantErr: false,
		},
		{
			name: "redirect",
			resp: AIAssistantResponse{Redirects: []AIAssistantRedirect{{
				Title: "Open Meshery Models",
				URL:   "/extensions/models",
			}}},
			wantErr: false,
		},
		{
			name: "structured error",
			resp: AIAssistantResponse{Errors: []AIAssistantError{{
				Code:    "CONTEXT_TOO_LARGE",
				Message: "Narrow the selected context.",
			}}},
			wantErr: false,
		},
		{
			name: "multiple output classes",
			resp: AIAssistantResponse{
				Explanation: "This design contains one service.",
				Recommendations: []AIAssistantRecommendation{{
					Title: "Review service selectors",
				}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.resp.Validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
