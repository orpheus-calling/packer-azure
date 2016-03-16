// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See the LICENSE file in the project root for license information.

package common

import (
	"fmt"
	"testing"
	"time"
)

func TestInterruptibleTaskShouldImmediatelyEndOnCancel(t *testing.T) {
	testSubject := NewInterruptibleTask(
		func() bool { return true },
		func() error {
			for {
				time.Sleep(time.Second * 30)
			}
		})

	result := testSubject.Run()
	if result.IsCancelled != true {
		t.Fatalf("Expected the task to be cancelled, but it was not.")
	}
}

func TestInterruptibleTaskShouldRunTaskUntilCompletion(t *testing.T) {
	var count int

	testSubject := &InterruptibleTask{
		IsCancelled: func() bool {
			return false
		},
		Task: func() error {
			for i := 0; i < 10; i++ {
				count += 1
			}

			return nil
		},
	}

	result := testSubject.Run()
	if result.IsCancelled != false {
		t.Errorf("Expected the task to *not* be cancelled, but it was.")
	}

	if count != 10 {
		t.Errorf("Expected the task to wait for completion, but it did not.")
	}

	if result.Err != nil {
		t.Errorf("Expected the task to return a nil error, but got=%s", result.Err)
	}
}

func TestInterruptibleTaskShouldImmediatelyStopOnTaskError(t *testing.T) {
	testSubject := &InterruptibleTask{
		IsCancelled: func() bool {
			return false
		},
		Task: func() error {
			return fmt.Errorf("boom")
		},
	}

	result := testSubject.Run()
	if result.IsCancelled != false {
		t.Errorf("Expected the task to *not* be cancelled, but it was.")
	}

	if result.Err == nil {
		t.Errorf("Expected the task to return an error, but it did not.")
	}
}
