module Choose exposing (..)
import Random
import Array
import String


choose : String -> Random.Generator Int
choose text =
    let 
        maxIndex = (List.length (String.indexes " " text)) - 1
    in
    Random.int 0 maxIndex

readWord text indexes index =
    let 
        start =
            case Array.get index indexes of 
            Just n ->
                n+1
            Nothing ->
                0
        stop =
            case  Array.get (index+1) indexes of
            Just n ->
                n
            Nothing ->
                0
    in   
        String.slice start stop text 
