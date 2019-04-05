# GOlang Betaface API client

### What is Betaface API?
[Betaface API](https://www.betafaceapi.com/wpa/) is a face detection and face recognition web service. Some interesting data returned by the API:
 - faces location within the image (+100 points per face)
 - gender
 - age
 - ethnicity
 - emotions

### GOlang Client
This betafaceapi package is a golang client that covers the Betaface API v2.

### Installation
```bash
go get github.com/enriquebris/gobetafaceapi
```

### Examples

##### Upload a new image and get all information related to: faces, faces attributes, faces recongnition
```go
package main

import (
	"fmt"

	"github.com/enriquebris/gobetafaceapi"
)

func main() {
	var (
		// free subscription plan credentials
		apiKey    = "d45fd466-51e2-4701-8da8-04351c872236"
		apiSecret = "no"
	)

	client := gobetafaceapi.NewNativeHTTPClient(apiKey, apiSecret)

	// upload an image
	httpResp, resp, errResp, err := client.SendMedia(
		"https://cdn.vox-cdn.com/thumbor/3V8wxIEwW8-JjMu-dX7lwcVQWd0=/0x0:1000x563/1200x800/filters:focal(420x202:580x362)/cdn.vox-cdn.com/uploads/chorus_image/image/60350569/killing_eve_review.0.jpg",
		gobetafaceapi.DetectionFlags{
			Classifiers: true,
		},
		[]string{},
		"killing_eve_review.0.jpg",
	)

	if err != nil {
		// request error
		fmt.Println(err)
		return
	}

	if httpResp.GetStatusCode() == 200 {
		// total faces
		fmt.Printf("total faces: %v\n", len(resp.Media.Faces))

		// print faces' tags
		for i := 0; i < len(resp.Media.Faces); i++ {
			fmt.Printf("\nface UUID: %v\n", resp.Media.Faces[i].FaceUUID)

			// print tags
			for c := 0; c < len(resp.Media.Faces[i].Tags); c++ {
				fmt.Printf("%v: %v [%v]\n", resp.Media.Faces[i].Tags[c].Name, resp.Media.Faces[i].Tags[c].Value, resp.Media.Faces[i].Tags[c].Confidence)
			}
		}
	} else {
		// the API returned an error
		fmt.Println(errResp.Description)
	}
}

```

### TODO

#### admin
- [ ] [GET] /v2/admin/usage

#### face
- [ ] [GET] /v2/face
- [ ] [GET] /v2/face/cropped

#### media
- [x] [GET] /v2/media
- [x] [POST] /v2/media
- [ ] [GET] /v2/media/hash
- [ ] [POST] /v2/media/file

#### person
- [ ] [GET] /v2/person
- [ ] [POST] /v2/person

#### recognize
- [ ] [GET] /v2/recognize
- [ ] [POST] /v2/recognize

#### transform
- [ ] [GET] /v2/transform
- [ ] [POST] /v2/transform