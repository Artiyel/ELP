module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)
import Html.Events exposing (onClick)
import Http
import Random
import Array
import Fetch_def exposing (getdef, Msg_2(..))
import Choose exposing (..)
import Compare exposing (..)



-- MAIN


main = Browser.element
        { init = init -- Initialisation du modèle
        , update = update -- Gestion des messages
        , subscriptions = subscriptions -- Pour Browser.element
        , view = view -- Affichage
        }
        

-- MODEL


type alias Model =
  { 
    title : String, -- Titre de l’application
    input : String, -- Saisie de l’utilisateur
    mot : String, -- Mot à deviner
    affichersolu : Bool, -- Indique si la solution est affichée
    totalmots : String, -- Texte contenant tous les mots
    def : List (String, List String),  -- Liste des définitions du mot
    reponse : Int  -- Indique si la réponse est correcte (0 ou 1)
  }

-- Initialisation du modèle et chargement du fichier texte

init : () -> ( Model, Cmd Msg )
init _ =
    ( { title = "Guess It!" 
      , input = "" 
      , mot = "" 
      , affichersolu = False 
      , def = []  
      , totalmots = "" 
      , reponse = 0
      }, 
      -- Requête HTTP pour charger le fichier contenant les mots
      Http.get 
        { url = "/static/thousand_words_things_explainer.txt"
        , expect = Http.expectString GotText 
        }
    )


-- UPDATE


type Msg = Change String  -- Mise à jour de la saisie utilisateur
        | Afficher -- Afficher la solution
        | FetchMsg Msg_2  -- Message provenant du module Fetch_def
        | GotText (Result Http.Error String) -- Résultat du chargement du fichier
        | Reponse Int  -- Résultat de la comparaison
        | NewIndex Int -- Index du mot choisi aléatoirement
        | NouveauMot -- Générer un nouveau mot

  

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
    
        -- Réception du fichier texte avec tous les mots
        GotText (Ok texte) ->
            ( { model | totalmots = texte }
            -- Génère un index aléatoire pour choisir un mot
            , Random.generate NewIndex (Choose.choose texte)
            )
   
        GotText (Err _) ->
            ( { model | totalmots = "Erreur de chargement" }, Cmd.none ) 
        
         -- Réception de l’index du mot choisi
        NewIndex index ->
            let 
                -- Création d’un tableau contenant les positions des espaces
                indexesArray = Array.fromList (String.indexes " " model.totalmots)
                -- Lecture du mot correspondant à l’index
                motChoisi = Choose.readWord model.totalmots indexesArray index
            in
            ( { model | mot = motChoisi }
            -- Récupération des définitions du mot
            , Cmd.map FetchMsg (getdef motChoisi)
            )

        -- Message provenant du module Fetch_def
        FetchMsg subMsg ->
                    case subMsg of
                        -- Définitions trouvées
                        Trouvé result -> 
                            case result of 
                                Ok listeDefinitions ->
                                    let
                                        -- Pour le mot récupéré, on prend chaque sens (Meaning)
                                        -- et on crée une paire : (nature grammaticale, liste de définitions)
                                        toutesLesDefs = 
                                            case List.head listeDefinitions of
                                                Just d -> 
                                                   List.map (\m -> (m.partOfSpeech, m.definitions)) d.meanings

                                                Nothing -> []
                                    in
                                    ( { model | def = toutesLesDefs }, Cmd.none )

                                -- Erreur lors de la récupération des définitions
                                Err _ ->
                                    ( { model | def = [] }, Cmd.none )

                        Rien -> -- Cas où le message de Fetch_def est 'Rien'
                            ( { model | def = [] }, Cmd.none )

        -- Mise à jour de la saisie utilisateur
        Change newContent ->
            ( { model | input = newContent }, Cmd.none )
        
        -- Affichage de la solution
        Afficher -> 
            ( { model | affichersolu = True }, Cmd.none )

        -- Vérification de la réponse
        Reponse reponse ->
            if Compare.isTheSame model.input model.mot == 1 then  
                ( { model | reponse = 1}, Cmd.none )
            else 
                ( { model | reponse = 0}, Cmd.none )
        
        -- Génération d’un nouveau mot et on nettoie 
        NouveauMot ->
            ( { model
                | input = ""
                , affichersolu = False
                , def = []
                , reponse = 0
            }
            -- Nouveau tirage aléatoire
            , Random.generate NewIndex (Choose.choose model.totalmots)
            )
          
    
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
       [ text model.title ],  -- Titre
  viewSolu model, -- Affichage de la solution
  div [style "margin-bottom" "40px", 
      style "margin-left" "100px"]
      (List.map afficherLigne model.def), -- Liste des définitions
  div [style "margin-bottom" "10px", 
       style "text-align" "center"] 
      [input [ placeholder "Try to guess the word...", value model.input, onInput Change ][]], -- Champ de saisie
  viewValidation model,  -- Message de validation
  div [style "margin-bottom" "20px", 
       style "text-align" "center"] 
      [button [ onClick Afficher ] [ text "Show the solution" ]], -- Bouton pour afficher la solution
  div [style "position" "fixed",
       style "bottom" "20px",
       style "right" "20px"] 
      [button [ onClick NouveauMot ] [ text "New word !" ]] -- Bouton pour générer un nouveau mot
  ]

-- Affichage du message de validation
viewValidation : Model -> Html Msg
viewValidation model =
    if model.input == "" then
        text "" 
    else if Compare.isTheSame model.input model.mot == 1 then 
        div [ style "color" "green", style "text-align" "center", style "margin-bottom" "10px" ] [ text "Congrats, you're right !" ]
    else 
        div [ style "color" "red", style "text-align" "center", style "margin-bottom" "10px"] [ text "This is not the right word..." ]

-- Affichage conditionnel de la solution
viewSolu : Model -> Html Msg
viewSolu model =
  if model.affichersolu == True then 
    div [ style "margin-left" "100px", style "font-size" "30px" ] [ text model.mot]

  else
    div [ style "color" "green" ] [ text ""]
  
-- Affichage d’une définition sous forme de liste
afficherLigne : (String, List String) -> Html Msg
afficherLigne (nat, defs) =
    div [ style "margin-bottom" "20px" ]
        -- Nature grammaticale : 
        [ h3 [] [ text nat ]  -- Transforme la chaîne de caractères en élément texte HTML et place ce texte dans une balise h3 (titre de niveau 3)
        -- Définitions : 
        , ul [] (List.map (\d -> li [] [ text d ]) defs)] -- Prend chaque définition d de la liste defs
                                                         -- et transforme chaque chaîne de caractères en élément texte (li) HTML
                                                         -- puis place tous ces éléments dans une liste (ul) HTML.
        