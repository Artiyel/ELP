module Fetch_def exposing (getdef, Msg_2(..), Definition)

import Http
import Json.Decode exposing (Decoder, map2, map3, map, field, string, list) 
 

-- Messages envoyés au module principal après la requête HTTP
type Msg_2 = Trouvé (Result Http.Error (List Definition)) -- Définitions récupérées ou erreur
            | Rien  -- Aucun résultat


-- Commande HTTP : Lance une requête vers l’API dictionnaire pour récupérer les définitions d’un mot
getdef : String -> Cmd Msg_2
getdef mot =
  Http.get
    { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ mot
    , expect = Http.expectJson Trouvé definitionsDecoder
    }

-- Structure représentant une entrée du dictionnaire
type alias Definition =
    { mot : String  -- Mot recherché
    , meanings : List Meaning -- Liste des sens du mot
    , prononciation : String -- Prononciation du mot
    }

-- Structure représentant un sens du mot
type alias Meaning =
    { partOfSpeech : String  -- Nature grammaticale (nom, verbe, etc.)
    , definitions : List String -- Liste des définitions textuelles
    }


-- DECODERS
-- Les decoders transforment le JSON de l’API en types Elm

-- Decoder pour une définition simple (texte)
definitionTextDecoder : Decoder String
definitionTextDecoder =
    field "definition" string

-- Decoder pour un "meaning"
meaningDecoder : Decoder Meaning
meaningDecoder =
    map2 Meaning
        (field "partOfSpeech" string)
        (field "definitions" (list definitionTextDecoder))

-- Decoder pour une entrée complète du dictionnaire
definitionDecoder : Decoder Definition
definitionDecoder =
    map3 Definition
        (field "word" string)
        (field "meanings" (list meaningDecoder) |> map (List.sortBy .partOfSpeech))
        (field "phonetic" string)


-- Decoder final : transforme la réponse JSON (liste d’entrées) en List Definition
definitionsDecoder : Decoder (List Definition)
definitionsDecoder =
    list definitionDecoder
