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
    [x]->[x]
    [a::b::xs]->
        if a==b then compress (b::xs)
        else compress a :: compress (b::xs)