import http from 'k6/http'
import {sleep} from 'k6'

export const options = {
    iterations: 10_000,   
    vus: 50,
    batch: 1_000,
}

export default function () {
    http.get("http://sample-service-custom.sample:5000/")
    sleep(.01)
}
