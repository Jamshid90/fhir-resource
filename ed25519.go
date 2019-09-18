package main

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"testing"
)

func main() {


	// hash : 15633aa88957225b1124473886d3200e3dc2842c02f6ac5bb6601b6fa767b074d6c36a0e09938e6daac0f85e51701a286af7b59b7b8d535c7fb086c96cdce986
	// public key : 37533e61ff0be2ad0e8552026f62d1123c871bfd02ed2758dcb1764a7527bd90
	// private key : 62788941319f86697d5e6010ed6111d27833ddf0c6157004612b6ca07f299fd937533e61ff0be2ad0e8552026f62d1123c871bfd02ed2758dcb1764a7527bd90

	//public, private, _ := ed25519.GenerateKey(nil)

	//public := []byte("37533e61ff0be2ad0e8552026f62d1123c871bfd02ed2758dcb1764a7527bd90")
	//private := []byte("62788941319f86697d5e6010ed6111d27833ddf0c6157004612b6ca07f299fd937533e61ff0be2ad0e8552026f62d1123c871bfd02ed2758dcb1764a7527bd90")


	var private ed25519.PrivateKey

	private = []byte("62788941319f86697d5e6010ed6111d27833ddf0c6157004612b6ca07f299f12")
	public := private.Public()


	data := `{
  "resourceType": "Appointment",
  "id": "example",
  "text": {
    "status": "generated",
    "div": "<div xmlns=\"http://www.w3.org/1999/xhtml\">Brian MRI results discussion</div>"
  },
  "status": "booked",
  "serviceCategory": [
    {
      "coding": [
        {
          "system": "http://example.org/service-category",
          "code": "gp",
          "display": "General Practice"
        }
      ]
    }
  ],
  "serviceType": [
    {
      "coding": [
        {
          "code": "52",
          "display": "General Discussion"
        }
      ]
    }
  ],
  "specialty": [
    {
      "coding": [
        {
          "system": "http://snomed.info/sct",
          "code": "394814009",
          "display": "General practice"
        }
      ]
    }
  ],
  "appointmentType": {
    "coding": [
      {
        "system": "http://terminology.hl7.org/CodeSystem/v2-0276",
        "code": "FOLLOWUP",
        "display": "A follow up visit from a previous appointment"
      }
    ]
  },
  "reasonReference": [
    {
      "reference": "Condition/example",
      "display": "Severe burn of left ear"
    }
  ],
  "priority": 5,
  "description": "Discussion on the results of your recent MRI",
  "start": "2013-12-10T09:00:00Z",
  "end": "2013-12-10T11:00:00Z",
  "created": "2013-10-10",
  "comment": "Further expand on the results of the MRI and determine the next actions that may be appropriate.",
  "basedOn": [
    {
      "reference": "ServiceRequest/myringotomy"
    }
  ],
  "participant": [
    {
      "actor": {
        "reference": "Patient/example",
        "display": "Peter James Chalmers"
      },
      "required": "required",
      "status": "accepted"
    },
    {
      "type": [
        {
          "coding": [
            {
              "system": "http://terminology.hl7.org/CodeSystem/v3-ParticipationType",
              "code": "ATND"
            }
          ]
        }
      ],
      "actor": {
        "reference": "Practitioner/example",
        "display": "Dr Adam Careful"
      },
      "required": "required",
      "status": "accepted"
    },
    {
      "actor": {
        "reference": "Location/1",
        "display": "South Wing, second floor"
      },
      "required": "required",
      "status": "accepted"
    }
  ]
}`;

	sha_512 := sha512.New()
	sha_512.Write([]byte(data))

	hmac512 := hmac.New(sha512.New, []byte("secret"))
	hmac512.Write([]byte(data))

	hash :=  fmt.Sprintf("%x", sha_512.Sum(nil))

	fmt.Println(hash)
	fmt.Println("----------------")


	sig := ed25519.Sign(private, []byte(hash))



	verify := ed25519.Verify(public, []byte(hash), sig)


	fmt.Println(sig)

	fmt.Println(verify)



}



func TestSignVerify(t *testing.T) {

	public, private, _ := ed25519.GenerateKey(nil)

	message := []byte("Test message")

	sig := ed25519.Sign(private, message)

	verify := ed25519.Verify(public, message, sig)

	fmt.Println(verify)

	if !verify {
		fmt.Println("valid signature rejected")
	}

	wrongMessage := []byte("wrong message")

	verify2 := ed25519.Verify(public, wrongMessage, sig)

	fmt.Println(verify2)

	if verify2{
		fmt.Println("signature of different message accepted")
	}
}
