package lib

import (
	"encoding/json"

	pb "github.com/rwrrioe/pythia_protos/gen/go/ocr"
)

func ConvertProto(resp *pb.OCRResponse) ([]byte, error) {
	jsonStr, err := json.Marshal(resp)
	if err != nil {
		return nil, nil
	}

	return jsonStr, nil
}
