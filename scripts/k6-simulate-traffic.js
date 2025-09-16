import http from 'k6/http'
import {sleep} from 'k6'

export const options = {
    iterations: 1000,   
    vus: 5,
    batch: 250,
}

export default function () {
    http.get("http://localhost:5000/")
    sleep(0.01)
}