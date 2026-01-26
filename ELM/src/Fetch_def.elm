module Fetch_def exposing (getdef, Msg_2(..), Definition)

import Http
import Json.Decode exposing (Decoder, map2, map3, map, field, string, list)

 -- Definitions de la structure de donnée utilisée
type Msg_2 = Trouvé (Result Http.Error (List Definition)) | Rien

type alias Definition =
    { mot : String
    , meanings : List Meaning
    , prononciation : String
    }


type alias Meaning =
    { partOfSpeech : String
    , definitions : List String
    }


-- COMMANDE HTTP
getdef : String -> Cmd Msg_2
getdef mot =
    Http.get
        { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ mot
        , expect = Http.expectJson Trouvé definitionsDecoder
        }


-- DECODERS dans "l'ordre" d'appel

definitionsDecoder : Decoder (List Definition)
definitionsDecoder =
    list definitionDecoder

definitionDecoder : Decoder Definition
definitionDecoder =
    map3 Definition
        (field "word" string)
        (field "meanings" (list meaningDecoder) |> map (List.sortBy .partOfSpeech))
        (field "phonetic" string)


meaningDecoder : Decoder Meaning
meaningDecoder =
    map2 Meaning
        (field "partOfSpeech" string)
        (field "definitions" (list definitionTextDecoder))

definitionTextDecoder : Decoder String
definitionTextDecoder =
    field "definition" string






