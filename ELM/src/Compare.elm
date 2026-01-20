module Compare exposing (..)

isTheSame : String -> String -> Int
isTheSame str1 str2 =
    if String.toLower str1 == String.toLower str2 then
        1
    else
        0

