module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)
import Html.Events exposing (onClick)
import Http



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
    totalmots : String
  }


init : () -> ( Model, Cmd Msg )
init _ =
    ( { title = "Guess It!"
      , input = ""
      , mot = ""
      , affichersolu = False
      , totalmots = ""
      }
    , Http.get 
        { url = "/static/thousand_words_things_explainer.txt"
        , expect = Http.expectString GotText 
        }
    )



-- UPDATE


type Msg = Change String | Afficher | GotText (Result Http.Error String)
  

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotText (Ok texte) ->
            ( { model | totalmots = texte }, Cmd.none )
        GotText (Err texte) ->
            ( { model | totalmots = "Erreur de chargement" }, Cmd.none )

        Change newContent ->
            ( { model | input = newContent }, Cmd.none )
          
        Afficher -> 
            ( { model | affichersolu = True }, Cmd.none )
        

    
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
        style "margin-left" "550px" ] 
       [ text model.title ],
  div [style "margin-bottom" "20px", 
       style "margin-left" "540px"] 
      [input [ placeholder "Try to guess the word...", value model.input, onInput Change ][]],
  div [style "margin-bottom" "20px", 
       style "margin-left" "560px"] 
      [button [ onClick Afficher ] [ text "Afficher la solution" ]]
  ]
  
  