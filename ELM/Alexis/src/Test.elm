module Test exposing (..)

addElemInList lenombre fois danslaliste =
    if fois == 0 then danslaliste
    else addElemInList lenombre (fois-1) (lenombre::danslaliste)

dupli liste =
    case liste of
    (x::xs) -> x::x:: dupli xs
    [] -> []

compress liste =
    case liste of
    [] -> []
    numero1::xs ->
        if List.length xs > 2
            then case xs of 
                numero2::xss -> if numero1==numero2 then compress numero2::xss else compress xss
                [a,b] ->  if numero1==numero2 then compress numero2::xss else compress xss
        else if numero1==xs then [xs]
        else if xs == [] then []
        else xs