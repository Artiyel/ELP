module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)
import Html.Events exposing (onClick)
import Http

import Fetch_def exposing (getdef, Msg_2(..))



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
    def : String
  }


init : () -> ( Model, Cmd Msg )
init _ =
    let
        motChoisi = "below"
    in
    ( { title = "Guess It!"
      , input = ""
      , mot = motChoisi
      , affichersolu = False
      , def = ""
      , totalmots = ""
      }, Cmd.map FetchMsg (getdef motChoisi) 
    )



-- UPDATE


type Msg = Change String | Afficher | FetchMsg Fetch_def.Msg_2 --| GotText (Result Http.Error String) |
  

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
    {-
        GotText (Ok texte) ->
            let 
                motChoisi = "below" 
            in
            ( { model | mot = motChoisi }
            , Cmd.map GotDef (getdef motChoisi)
            )
        
        GotText (Err error) ->
            ( { model | totalmots = "Erreur de chargement du fichier" }, Cmd.none )
      -}

        FetchMsg result ->
                    case result of
                        Ok listeDefinitions ->
                            let 
                                premierMot = List.head listeDefinitions 
                                texteAAfficher = 
                                    case premierMot of
                                        Just d -> String.join ", " d.definitions
                                        Nothing -> "Aucune définition trouvée"
                            in
                            ( { model | def = texteAAfficher }, Cmd.none )

                        Err _ ->
                            ( { model | def = "Erreur HTTP" }, Cmd.none )


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
       style "margin-left" "560px"] 
      [ text model.def],
  div [style "margin-bottom" "20px", 
       style "margin-left" "540px"] 
      [input [ placeholder "Try to guess the word...", value model.input, onInput Change ][]],
  div [style "margin-bottom" "20px", 
       style "margin-left" "560px"] 
      [button [ onClick Afficher ] [ text "Afficher la solution" ]]
  ]
  
  