module Compare exposing (..)

-- Fonction permettant de comparer deux chaînes de caractères
-- La comparaison est insensible à la casse (majuscules / minuscules)

isTheSame : String -> String -> Int
isTheSame str1 str2 =
    if String.toLower str1 == String.toLower str2 then
        1
    else
        0

