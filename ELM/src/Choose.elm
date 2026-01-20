module Choose exposing (..)
import Random
import Array
import String

test text =
    let seed = Random.initialSeed 12
    in choose text seed

choose text seed =
    let 
        ( index, nextSeed ) =
            Random.step (Random.int 0 1000) seed
    in
        (readWord text (Array.fromList (String.indexes " " text)) index , nextSeed)

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
