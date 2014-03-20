package webapi

import (
        "bytes"
        "encoding/json"
        "encoding/xml"
        "fmt"
)

type WebApiEncoder func(v interface{}) ([]byte, error);

func JsonEncoder(v interface{}) ([]byte, error) {
    b, err := json.Marshal(v)
    return b, err
}

func XmlEncoder(v interface{}) ([]byte, error) {
    var buf bytes.Buffer
    
    if _, err := buf.Write([]byte(xml.Header)); err != nil {
        return nil, err
    }
    if _, err := buf.Write([]byte("<root>")); err != nil {
        return nil, err
    }
    b, err := xml.Marshal(v)

    if err != nil {
        return nil, err
    }
    if _, err := buf.Write(b); err != nil {
        return nil, err
    }
    if _, err := buf.Write([]byte("</root>")); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func TextEncoder(v interface{}) ([]byte, error) {
    
    var buf bytes.Buffer

    // for _, v := range v {
        if _, err := fmt.Fprintf(&buf, "%s\n", v); err != nil {
            return nil, err
        }
    // }
    
    return buf.Bytes(), nil
}