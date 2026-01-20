module Fetch_def exposing (..)

import Http
import Json.Decode exposing (Decoder, map3, map, field, string, list)

type Msg = Trouvé (Result Http.Error (List Definition)) | Rien

type alias Definition =
  { mot : String
  , definitions : List String
  , prononciation : String
  }


getdef : String -> Cmd Msg
getdef mot =
  Http.get
    { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ mot
    , expect = Http.expectJson Trouvé definitionsDecoder
    }

definitionTextDecoder : Decoder String
definitionTextDecoder =
  field "definition" string


meaningDefinitionsDecoder : Decoder (List String)
meaningDefinitionsDecoder =
  field "definitions" (list definitionTextDecoder)


allDefinitionsDecoder : Decoder (List String)
allDefinitionsDecoder =
  field "meanings" (list meaningDefinitionsDecoder)
    |> map List.concat


definitionDecoder : Decoder Definition
definitionDecoder =
  map3 Definition
    (field "word" string)
    allDefinitionsDecoder
    (field "phonetic" string)


definitionsDecoder : Decoder (List Definition)
definitionsDecoder =
  list definitionDecoder
