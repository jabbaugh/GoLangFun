package main

import "bufio"
import "os"
import "net/http"
import "fmt"
import "io"
import "io/ioutil"
import "encoding/json"
import "strings"
import "time"

type PartnerMappings []struct {
    ID                   string `json:"id"`
    FulfillmentChannelID int    `json:"fulfillmentChannelId"`
    AudienceID           int    `json:"audienceId"`
    PartnerAccountID     string `json:"partnerAccountId"`
    ChannelInputs        struct {
      ChildPartnerID string `json:"childPartnerId"`
      Subtype        string `json:"subtype"`
    } `json:"channelInputs"`
    PartnerAudienceID string `json:"partnerAudienceId"`
    ChannelOutputs    struct {
      ChildPartnerID   string `json:"childPartnerId"`
      PartnerAccountID string `json:"partnerAccountId"`
    } `json:"channelOutputs"`
    Active     bool      `json:"active"`
    Source     string    `json:"source"`
    CreatedAt  time.Time `json:"createdAt"`
    ModifiedAt time.Time `json:"modifiedAt"`
  } 

type Response struct {
  PageNo          int `json:"pageNo"`
  PageSize        int `json:"pageSize"`
  TotalPages      int `json:"totalPages"`
  TotalResults    int `json:"totalResults"`
  PartnerMappings PartnerMappings `json:"partnerMappings"`
}

func main() {
//  audId := "3785834"

  readFile, err := os.Open("audiences-uniq.txt")

  if err != nil {
    fmt.Println(err)
  }

  fileScanner := bufio.NewScanner(readFile)

  fileScanner.Split(bufio.ScanLines)

  for fileScanner.Scan() {
    partnerMappings := GetPartnerMappings(fileScanner.Text())

    deleteRes := ""
    for _, rec := range partnerMappings {
      deleteRes = DeletePartnerMapping(rec.ID)
      fmt.Println(deleteRes)
    }

    if ! strings.Contains(deleteRes, "NOT_FOUND") {
      time.Sleep(1 * time.Second)
    }
  }

  readFile.Close()

}

func GetPartnerMappings(audId string) PartnerMappings {
  fmt.Println(audId)
  resp, err := http.Get("https://" + audId)
  if err != nil {
      // handle error
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
      // handle error
  }

  var result Response
  if err := json.Unmarshal(body, &result); err != nil {   // Parse []byte to go struct pointer
      fmt.Println("Can not unmarshal JSON")
  }

  return result.PartnerMappings
}

func DeletePartnerMapping(id string) string {
  fmt.Println("Delete: " + id)

    // create a new HTTP client
  client := &http.Client{}

  // create a new DELETE request
  req, err := http.NewRequest("DELETE", "https://" + id, nil)
  if err != nil {
      panic(err)
  }

  // send the request
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  // read the response body
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      panic(err)
  }


  return string(body)
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}
