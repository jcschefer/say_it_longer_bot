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
   "os"
   "strings"
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
   word := "map"
   fmt.Println( get_longest_synonym( word, thesaurus_key ) )
   s := api.PublicStreamSample(nil)
   data := <-s.C
   for id := range data.([]int) {
      tw, err := api.GetTweet( int64(id), nil )
      check( err )
      fmt.Println(tw.Text)
   }
   //
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
   anaconda.SetConsumerKey( keys[ 0 ] )
   anaconda.SetConsumerSecret( keys[ 1 ] )
   a := anaconda.NewTwitterApi( keys[ 2 ], keys[ 3 ] )
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
   longest := ""
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
