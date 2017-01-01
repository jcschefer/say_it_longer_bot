// Jack Schefer, began 12/10/16
//
package main
//
import(
   "bufio"
   "encoding/xml"
   "fmt"
   "io/ioutil"
   "net/http"
   "net/url"
   "os"
   "strings"
   "strconv"
   //
   "github.com/ChimeraCoder/anaconda"
)
//
type Sense struct{
   Synonyms string   `xml:"syn"`
}
//
type Entry struct {
   Senses   []Sense  `xml:"sens"`
}
//
type Result struct {
   XMLName  xml.Name `xml:"entry_list"`
   Entries  []Entry  `xml:"entry"`
}
//
////////////////////////////////////////////////////////////////////
//
func main() {
   //
   api, thesaurus_key := setup_api()
   defer api.Close()
   //
   self, err := api.GetSelf(nil)
   check(err)
   //
   phrase := "Hi my name is jack, welcome to my house"
   fmt.Println(phrase)
   fmt.Println(longer_phrase(phrase, thesaurus_key))
   //
   word := "map"
   fmt.Println( get_longest_synonym( word, thesaurus_key ) )
   vals := url.Values{}
   vals.Set("follow", strconv.FormatInt(self.Id, 10))
   s := api.PublicStreamFilter(vals)
   fmt.Println("Listening...")
   for {
      data := <-s.C
      switch v:= data.(type) {
         default:
            fmt.Printf("unexpected type %T\n", v)
         case anaconda.Tweet:
            fmt.Println(data.(anaconda.Tweet).Text)
         case anaconda.StatusDeletionNotice:
            fmt.Println("Status Deletion Notice")
      }
   }
   //
}
//
////////////////////////////////////////////////////////////////////
//
func longer_phrase( s string, key string ) string {
   words := strings.Split(s, " ")
   for i,_ := range words {
      words[ i ] = strings.Trim( words[ i ], " \t\n!.,:;()[]{}@#$%^&*-+" )
      words[ i ] = get_longest_synonym( words[ i ], key )
   }
   st := ""
   for _,w := range words {
      st += w
      st += " "
   }
   return st
}
//
////////////////////////////////////////////////////////////////////
//
func setup_api() ( *anaconda.TwitterApi, string ) {
   //
   file, err := os.Open( "keys.txt" )
   check(err)
   defer file.Close()
   //
   var keys [ 5 ]string
   scanner := bufio.NewScanner( file )
   i := 0
   //
   for scanner.Scan() {
     keys[ i ] = scanner.Text()
     i += 1
   }
   //
   anaconda.SetConsumerKey(      keys[ 0 ]            )
   anaconda.SetConsumerSecret(   keys[ 1 ]            )
   a := anaconda.NewTwitterApi(  keys[ 2 ], keys[ 3 ] )
   return a, keys[ 4 ]
   //
}
//
////////////////////////////////////////////////////////////////////
//
func get_longest_synonym( word string, key string) string {
   //
   response, err := http.Get( "http://www.dictionaryapi.com/api/v1/references/thesaurus/xml/" + word + "?key=" + key)
   check(err)
   defer response.Body.Close()
   body, err := ioutil.ReadAll( response.Body )
   //
   result := Result{}
   xml.Unmarshal( body, &result )
   //
   longest := word
   for _, entry := range result.Entries {
      for _, sense := range entry.Senses {
         syn_slice := strings.Split( sense.Synonyms, "," )
         for _, w := range( syn_slice ) {
            ind := strings.Index( w, "(" )
            if ind == -1 { ind = len(w) }
            trimmed := strings.TrimSpace( string( w[:ind ] ) )
            if len( trimmed ) > len( longest ) {
               longest = trimmed
            }
         }
      }
   }
   return longest
   //
}
//
////////////////////////////////////////////////////////////////////
//
func check( e interface{} ) {
   if e != nil {
      panic(e)
   }
}
//
////////////////////////////////////////////////////////////////////
//
// End of file.
