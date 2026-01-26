module Choose exposing (..)
import Random
import Array
import String

-- Module permettant de choisir un mot aléatoirement dans un texte

-- Génère un index aléatoire correspondant à un mot du texte, le texte est supposé être une suite de mots séparés par des espaces
choose : String -> Random.Generator Int
choose text =
    let 
        -- Calcul du nombre d’espaces dans le texte = nombre de mots - 1
        maxIndex = (List.length (String.indexes " " text)) - 1
    in
    Random.int 0 maxIndex

-- Extrait un mot du texte à partir d’un index donné
-- indexes contient les positions des espaces dans le texte
readWord text indexes index =
    let 
        -- La position de début du mot correspond au caractère juste après l’espace précédent
        start =
            case Array.get index indexes of 
            Just n ->
                n+1
            -- Si aucun espace n’est trouvé (début du texte)
            Nothing ->
                0
        -- La position de fin du mot correspond à la position de l’espace suivant
        stop =
            case  Array.get (index+1) indexes of
            Just n ->
                n
            Nothing ->
                0
    in   -- Extraction du mot entre start et stop
        String.slice start stop text 
