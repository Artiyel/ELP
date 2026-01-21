module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)
import Html.Events exposing (onClick)
import Http
import Random
import Fetch_def exposing (getdef, Msg_2(..))
import Choose exposing (..)
import Compare exposing (..)



-- MAIN


main = Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        }
        

-- MODEL


type alias Model =
  { 
    title : String,
    input : String,
    mot : String,
    affichersolu : Bool,
    totalmots : String,
    def : List String,
    seed : Random.Seed,
    reponse : Int
  }


init : () -> ( Model, Cmd Msg )
init _ =
    ( { title = "Guess It!"
      , input = ""
      , mot = ""
      , affichersolu = False
      , def = []
      , totalmots = ""
      , seed = Random.initialSeed 12
      , reponse = 0
      }, Http.get 
        { url = "/static/thousand_words_things_explainer.txt"
        , expect = Http.expectString GotText 
        }
    )


-- UPDATE


type Msg = Change String | Afficher | FetchMsg Msg_2 | GotText (Result Http.Error String) | Reponse Int
  

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
    
        GotText (Ok texte) ->
            let 
                ( motChoisi, nextSeed ) = Choose.choose texte model.seed
            in
            ( { model | mot = motChoisi, seed = nextSeed }
            , Cmd.map FetchMsg (getdef motChoisi) 
            )

        GotText (Err _) ->
            ( { model | totalmots = "Erreur de chargement" }, Cmd.none )

        FetchMsg subMsg ->
                    case subMsg of
                        Trouvé result -> 
                            case result of 
                                Ok listeDefinitions ->
                                    let
                                        toutesLesDefs = 
                                            case List.head listeDefinitions of
                                                Just d -> 
                                                    -- On récupère la liste de définitions de chaque "meaning" 
                                                    -- et on les fusionne en une seule grande liste
                                                    List.concatMap .definitions d.meanings

                                                Nothing -> []
                                    in
                                    ( { model | def = toutesLesDefs }, Cmd.none )

                                Err _ ->
                                    ( { model | def = [] }, Cmd.none )
                        Rien -> -- Cas où le message de Fetch_def est 'Rien'
                            ( { model | def = [] }, Cmd.none )


        Change newContent ->
            ( { model | input = newContent }, Cmd.none )
          
        Afficher -> 
            ( { model | affichersolu = True }, Cmd.none )

        Reponse reponse ->
            if Compare.isTheSame model.input model.mot == 1 then  
                ( { model | reponse = 1}, Cmd.none )
            else 
                ( { model | reponse = 0}, Cmd.none )
          
    
-- SUB 

subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none


-- VIEW


view : Model -> Html Msg
view model =
  div[]
  [div [style "margin-bottom" "40px", 
        style "font-size" "40px", 
        style "text-align" "center" ] 
       [ text model.title ],
  viewSolu model, 
  div [style "margin-bottom" "40px", 
       style "margin-left" "100px"] 
      [ ul [] (List.map afficherLigne model.def)],
  div [style "margin-bottom" "20px", 
       style "text-align" "center"] 
      [input [ placeholder "Try to guess the word...", value model.input, onInput Change ][]],
  div [style "margin-bottom" "20px", 
       style "text-align" "center"] 
      [button [ onClick Afficher ] [ text "Afficher la solution" ]],
  viewValidation model
  ]

viewValidation : Model -> Html Msg
viewValidation model =
    if model.input == "" then
        text "" 
    else if Compare.isTheSame model.input model.mot == 1 then 
        div [ style "color" "green", style "text-align" "center" ] [ text "Bravo, vous avez trouvé le mot !" ]
    else 
        div [ style "color" "red", style "text-align" "center" ] [ text "Ce n'est pas le bon mot..." ]


viewSolu : Model -> Html Msg
viewSolu model =
  if model.affichersolu == True then 
    div [ style "margin-left" "100px", style "font-size" "30px" ] [ text model.mot]

  else
    div [ style "color" "green" ] [ text ""]
  
  
afficherLigne : String -> Html Msg
afficherLigne texteDef =
    li [ style "margin-bottom" "10px" ] [ text texteDef ]