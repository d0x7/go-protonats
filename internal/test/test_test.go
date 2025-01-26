package test

import (
	"regexp"
	"slices"
	"testing"
	go_nats "xiam.li/go-nats"
)

func TestNormal(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for range 3 {
		id := NewTestServiceNATSServer(instance.Conn, &testImplementation{}, go_nats.WithoutLeaderFns(), go_nats.WithoutFollowerFns()).Info().ID
		ids = append(ids, id)
	}
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("TestTest", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.NormalTestTest(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		re := regexp.MustCompile(`^server replying to Test Client from ([a-zA-Z0-9]+)$`)
		matches := re.FindStringSubmatch(resp.Test)
		if len(matches) != 2 {
			t.Fatalf("Response format doesn't match: %v", resp.Test)
		}
		if !slices.Contains(ids, matches[1]) {
			t.Fatalf("Server ID not found in the list of IDs: %v", matches[1])
		}
	})

	t.Run("EmptyTest", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.NormalEmptyTest()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		re := regexp.MustCompile(`^server replying to empty from ([a-zA-Z0-9]+)$`)
		matches := re.FindStringSubmatch(resp.Test)
		if len(matches) != 2 {
			t.Fatalf("Response format doesn't match: %v", resp.Test)
		}
		if !slices.Contains(ids, matches[1]) {
			t.Fatalf("Server ID not found in the list of IDs: %v", matches[1])
		}
	})

	t.Run("TestEmpty", func(t *testing.T) {
		t.Parallel()
		err := cli.NormalTestEmpty(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
	})

	t.Run("EmptyEmpty", func(t *testing.T) {
		t.Parallel()
		err := cli.NormalEmptyEmpty()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
	})
}

func TestNormalBroadcast(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for range 3 {
		id := NewTestServiceNATSServer(instance.Conn, &testImplementation{}, go_nats.WithoutLeaderFns(), go_nats.WithoutFollowerFns()).Info().ID
		ids = append(ids, id)
	}
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("TestTest", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.NormalBroadcastTestTest(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
		if len(resp) != 3 {
			t.Fatalf("Expected 3 responses, got %d: %v", len(resp), resp)
		}
		for _, r := range resp {
			re := regexp.MustCompile(`^server replying to Test Client from ([a-zA-Z0-9]+)$`)
			matches := re.FindStringSubmatch(r.Test)
			if len(matches) != 2 {
				t.Fatalf("Response format doesn't match: %v", r.Test)
			}
			if !slices.Contains(ids, matches[1]) {
				t.Fatalf("Server ID not found in the list of IDs: %v", matches[1])
			}
		}
	})

	t.Run("EmptyTest", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.NormalBroadcastEmptyTest()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
		if len(resp) != 3 {
			t.Fatalf("Expected 3 responses, got %d: %v", len(resp), resp)
		}
		for _, r := range resp {
			re := regexp.MustCompile(`^server replying to empty from ([a-zA-Z0-9]+)$`)
			matches := re.FindStringSubmatch(r.Test)
			if len(matches) != 2 {
				t.Fatalf("Response format doesn't match: %v", r.Test)
			}
			if !slices.Contains(ids, matches[1]) {
				t.Fatalf("Server ID not found in the list of IDs: %v", matches[1])
			}
		}
	})

	t.Run("TestEmpty", func(t *testing.T) {
		t.Parallel()
		srvErrs, err := cli.NormalBroadcastTestEmpty(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
	})

	t.Run("EmptyEmpty", func(t *testing.T) {
		t.Parallel()
		srvErrs, err := cli.NormalBroadcastEmptyEmpty()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
	})
}

func TestErr(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	NewTestServiceNATSServer(instance.Conn, &testImplementation{}, go_nats.WithoutLeaderFns(), go_nats.WithoutFollowerFns())
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.ErrServiceError(&Test{Test: "Test Client"})
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if resp != nil {
			t.Fatalf("Expected nil response, got: %v", resp)
		}
		if !go_nats.IsServiceError(err) {
			t.Fatalf("Expected service error, got: %v", err)
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.ErrServerError(&Test{Test: "Test Client"})
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if resp != nil {
			t.Fatalf("Expected nil response, got: %v", resp)
		}
		if !go_nats.IsServiceError(err) {
			t.Fatalf("Expected service error, got: %v", err)
		}
	})
}

func TestErrBroadcast(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for range 3 {
		id := NewTestServiceNATSServer(instance.Conn, &testImplementation{}).Info().ID
		ids = append(ids, id)
	}
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.ErrServiceErrorBroadcast(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(resp) != 0 {
			t.Fatalf("Unexpected responses: %v", resp)
		}
		if len(srvErrs) != 3 {
			t.Fatalf("Expected 3 service errors, got %d: %v", len(srvErrs), srvErrs)
		}
		for _, e := range srvErrs {
			if !go_nats.IsServiceError(e) {
				t.Fatalf("Expected service error, got: %v", e)
			}
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.ErrServerErrorBroadcast(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(resp) != 0 {
			t.Fatalf("Unexpected responses: %v", resp)
		}
		if len(srvErrs) != 3 {
			t.Fatalf("Expected 3 service errors, got %d: %v", len(srvErrs), srvErrs)
		}
		for _, e := range srvErrs {
			if !go_nats.IsServiceError(e) {
				t.Fatalf("Expected service error, got: %v", e)
			}
		}
	})
}

func TestLeaderOnly(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	for range 3 {
		NewTestServiceNATSFollowerServer(instance.Conn, &testImplementation{})
	}
	leaderImpl := &testImplementation{}
	NewTestServiceNATSLeaderServer(instance.Conn, leaderImpl)
	if leaderImpl.id == "" {
		t.Fatalf("Server id is empty")
	}

	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("TestTest", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.LeaderOnlyTestTest(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if resp.Test != "leader replying to Test Client from "+leaderImpl.id {
			t.Fatalf("Unexpected response: %v", resp.Test)
		}
	})

	t.Run("EmptyTest", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.LeaderOnlyEmptyTest()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if resp.Test != "leader replying to empty from "+leaderImpl.id {
			t.Fatalf("Unexpected response: %v", resp.Test)
		}
	})

	t.Run("TestEmpty", func(t *testing.T) {
		t.Parallel()
		err := cli.LeaderOnlyTestEmpty(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
	})

	t.Run("EmptyEmpty", func(t *testing.T) {
		t.Parallel()
		err := cli.LeaderOnlyEmptyEmpty()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
	})
}

func TestLeaderOnlyBroadcast(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	for range 3 {
		NewTestServiceNATSFollowerServer(instance.Conn, &testImplementation{})
	}
	id := NewTestServiceNATSLeaderServer(instance.Conn, &testImplementation{}).Info().ID
	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("TestTest", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.LeaderOnlyBroadcastTestTest(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
		if len(resp) != 1 {
			t.Fatalf("Expected only one response, got %d: %v", len(resp), resp)
		}
		if resp[0].Test != "leader replying to Test Client from "+id {
			t.Fatalf("Unexpected response: %v", resp[0].Test)
		}
	})

	t.Run("EmptyTest", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.LeaderOnlyBroadcastEmptyTest()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
		if len(resp) != 1 {
			t.Fatalf("Expected only one response, got %d: %v", len(resp), resp)
		}
		if resp[0].Test != "leader replying to empty from "+id {
			t.Fatalf("Unexpected response: %v", resp[0].Test)
		}
	})

	t.Run("TestEmpty", func(t *testing.T) {
		t.Parallel()
		srvErrs, err := cli.LeaderOnlyBroadcastTestEmpty(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
	})

	t.Run("EmptyEmpty", func(t *testing.T) {
		t.Parallel()
		srvErrs, err := cli.LeaderOnlyBroadcastEmptyEmpty()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
	})
}

func TestFollowerOnly(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for range 3 {
		id := NewTestServiceNATSFollowerServer(instance.Conn, &testImplementation{}).Info().ID
		ids = append(ids, id)
	}
	NewTestServiceNATSLeaderServer(instance.Conn, &testImplementation{})

	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("TestTest", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.FollowerOnlyTestTest(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		re := regexp.MustCompile(`^follower replying to Test Client from ([a-zA-Z0-9]+)$`)
		matches := re.FindStringSubmatch(resp.Test)
		if len(matches) != 2 {
			t.Fatalf("Response format doesn't match: %v", resp.Test)
		}
		if !slices.Contains(ids, matches[1]) {
			t.Fatalf("Server ID not found in the list of follower IDs: %v", matches[1])
		}
	})

	t.Run("EmptyTest", func(t *testing.T) {
		t.Parallel()
		resp, err := cli.FollowerOnlyEmptyTest()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		re := regexp.MustCompile(`^follower replying to empty from ([a-zA-Z0-9]+)$`)
		matches := re.FindStringSubmatch(resp.Test)
		if len(matches) != 2 {
			t.Fatalf("Response format doesn't match: %v", resp.Test)
		}
		if !slices.Contains(ids, matches[1]) {
			t.Fatalf("Server ID not found in the list of follower IDs: %v", matches[1])
		}
	})

	t.Run("TestEmpty", func(t *testing.T) {
		t.Parallel()
		err := cli.FollowerOnlyTestEmpty(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
	})

	t.Run("EmptyEmpty", func(t *testing.T) {
		t.Parallel()
		err := cli.FollowerOnlyEmptyEmpty()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
	})
}

func TestFollowerOnlyBroadcast(t *testing.T) {
	t.Parallel()
	instance := newNATS(t)
	t.Cleanup(instance.Stop)
	var ids []string
	for range 3 {
		id := NewTestServiceNATSFollowerServer(instance.Conn, &testImplementation{}).Info().ID
		ids = append(ids, id)
	}
	NewTestServiceNATSLeaderServer(instance.Conn, &testImplementation{})

	cli := NewTestServiceNATSClient(instance.Conn)

	t.Run("TestTest", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.FollowerOnlyBroadcastTestTest(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
		if len(resp) != 3 {
			t.Fatalf("Expected 3 responses, got %d: %v", len(resp), resp)
		}
		for _, r := range resp {
			re := regexp.MustCompile(`^follower replying to Test Client from ([a-zA-Z0-9]+)$`)
			matches := re.FindStringSubmatch(r.Test)
			if len(matches) != 2 {
				t.Fatalf("Response format doesn't match: %v", r.Test)
			}
			if !slices.Contains(ids, matches[1]) {
				t.Fatalf("Server ID not found in the list of follower IDs: %v", matches[1])
			}
		}
	})

	t.Run("EmptyTest", func(t *testing.T) {
		t.Parallel()
		resp, srvErrs, err := cli.FollowerOnlyBroadcastEmptyTest()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
		if len(resp) != 3 {
			t.Fatalf("Expected 3 responses, got %d: %v", len(resp), resp)
		}
		for _, r := range resp {
			re := regexp.MustCompile(`^follower replying to empty from ([a-zA-Z0-9]+)$`)
			matches := re.FindStringSubmatch(r.Test)
			if len(matches) != 2 {
				t.Fatalf("Response format doesn't match: %v", r.Test)
			}
			if !slices.Contains(ids, matches[1]) {
				t.Fatalf("Server ID not found in the list of follower IDs: %v", matches[1])
			}
		}
	})

	t.Run("TestEmpty", func(t *testing.T) {
		t.Parallel()
		srvErrs, err := cli.FollowerOnlyBroadcastTestEmpty(&Test{Test: "Test Client"})
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
	})

	t.Run("EmptyEmpty", func(t *testing.T) {
		t.Parallel()
		srvErrs, err := cli.FollowerOnlyBroadcastEmptyEmpty()
		if err != nil {
			t.Fatalf("Error calling method: %v", err)
		}
		if len(srvErrs) != 0 {
			t.Fatalf("Unexpected server errors: %v", srvErrs)
		}
	})
}
