package cachemanager

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	rankManager "github.com/Ekram-B2/rankmanager/rank"
)

type testClientHit struct{}

func (t testClientHit) Get(key string) (interface{}, bool) {

	bytes, _ := json.Marshal(rankManager.Rank{Name: "Toronto", Rank: 0.5})
	return bytes, true
}

func (t testClientHit) Set(k string, x interface{}, d time.Duration) {
	return
}

type testClientMiss struct{}

func (t testClientMiss) Get(k string) (interface{}, bool) {
	return nil, false
}

func (t testClientMiss) Set(k string, x interface{}, d time.Duration) {
	return
}

func isSameSlice(one, two []byte) bool {
	for index, byteItem := range one {
		if byteItem != two[index] {
			return false
		}
	}
	return true
}

func Test_cachemanager_getBytesFromCache(t *testing.T) {
	// 1. Set up temp to store actual value of the rmclient (Arrange)
	tempclient := rmclient
	rmclient = testClientHit{}
	// 2. Define the expected and input args (Arrange)
	expectedFound := true
	expectedBytes, _ := json.Marshal(rankManager.Rank{Name: "Toronto", Rank: 0.5})
	inputKey := "hello"
	// 3. Define instance of ram cache manager (Arrange)
	rm := ramCacheManager{}
	// 4. Compute operation and get output (Act)
	actualIsFound, actualBytes := rm.getBytesFromCache(inputKey)
	// 5. Check to see if expected matches actual (Assert)
	if len(expectedBytes) != len(actualBytes) {
		t.Fatalf("expected did not match actual; expected was %v but actual was %v", expectedBytes, actualBytes)
	}

	if !isSameSlice(expectedBytes, actualBytes) {
		t.Fatalf("expected did not match actual; expected was %v but actual was %v", expectedBytes, actualBytes)
	}

	if expectedFound != actualIsFound {
		t.Fatalf("expected did not match actual; expected was %v but actual was %v", expectedFound, actualIsFound)
	}
	// 6. Reset the client
	rmclient = tempclient
}

func Test_cachemanager_putInCache(t *testing.T) {
	// 1. Set up temp to store actual value of the rmclient (Arrange)
	tempclient := rmclient
	rmclient = testClientHit{}
	// 2. Define input args (Arrange)
	inputEncoder := encodeRaw
	inputKey := "hello"
	// 3. Define instance of ram cache manager (Arrange)
	rm := ramCacheManager{}
	inputRank := rankManager.Rank{}
	// 4. Apply input args (Act)
	actualOut := rm.putInCache(inputKey, inputRank, inputEncoder)
	// 5. Check to see if expected matches actual (Assert)
	if actualOut != nil {
		t.Fatalf("Expected nil but got %v", actualOut)
	}
	// 6. Reset the client
	rmclient = tempclient

}

func Test_cachemanager_getRankFromCache(t *testing.T) {

	tests := []struct {
		name             string
		key              string
		decoderType      byteDecoder
		wantIsCacheHit   bool
		wantReturnedRank rankManager.Rank
		wantErr          bool
		client           ramCacheclient
	}{
		// 1. Setup expected and input to perform operation and then compare to see if
		// the actual matches what was expected (Arrange)
		{
			name:             "isNotHit",
			key:              "hello",
			decoderType:      decodeRaw,
			wantIsCacheHit:   false,
			wantReturnedRank: rankManager.Rank{},
			client:           testClientMiss{},
		},
		{
			name:             "isHit",
			key:              "hello",
			decoderType:      decodeRaw,
			wantIsCacheHit:   true,
			wantReturnedRank: rankManager.Rank{Name: "Toronto", Rank: 0.5},
			client:           testClientHit{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := ramCacheManager{}
			// 2. Set up temp to store actual value of the rmclient (Arrange)
			testClient := rmclient
			rmclient = tt.client
			// 3. Apply input to perform operation (Act)
			gotIsCacheHit, gotReturnedRank, err := rm.getRankFromCache(tt.key, tt.decoderType)
			// 4. Check to see if expected matches actual (Assert)
			if err != nil {
				t.Errorf("ramCacheManager.getRankFromCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsCacheHit != tt.wantIsCacheHit {
				t.Errorf("ramCacheManager.getRankFromCache() gotIsCacheHit = %v, want %v", gotIsCacheHit, tt.wantIsCacheHit)
			}
			if !reflect.DeepEqual(gotReturnedRank, tt.wantReturnedRank) {
				t.Errorf("ramCacheManager.getRankFromCache() gotReturnedRank = %v, want %v", gotReturnedRank, tt.wantReturnedRank)
			}
			// 5. Reset the client
			rmclient = testClient
		})
	}
}
