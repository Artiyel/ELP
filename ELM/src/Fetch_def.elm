module Fetch_def exposing (getDef)

import Http
import Json.Decode exposing (Decoder, map3, map, field, string, list)

type Msg_2 = Trouvé (Result Http.Error (List Definition)) | Rien

type alias Definition =
  { mot : String
  , definitions : List String
  , prononciation : String
  }


getdef : String -> Cmd Msg_2
getdef mot =
  Http.get
    { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ mot
    , expect = Http.expectJson Trouvé definitionsDecoder
    }

type alias Definition =
    { mot : String
    , meanings : List Meaning
    , prononciation : String
    }


-- COMMANDE HTTP
getDef : String -> Cmd Msg
getDef mot =
    Http.get
        { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ mot
        , expect = Http.expectJson Trouvé definitionsDecoder
        }


-- DECODERS
definitionTextDecoder : Decoder String
definitionTextDecoder =
    field "definition" string


meaningDecoder : Decoder Meaning
meaningDecoder =
    map2 Meaning
        (field "partOfSpeech" string)
        (field "definitions" (list definitionTextDecoder))


definitionDecoder : Decoder Definition
definitionDecoder =
    map3 Definition
        (field "word" string)
        (field "meanings" (list meaningDecoder) |> map (List.sortBy .partOfSpeech))
        (field "phonetic" string)


definitionsDecoder : Decoder (List Definition)
definitionsDecoder =
    list definitionDecoder
