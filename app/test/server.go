package main

import (
	"fmt"
	"net/http"
)

func soapHandler(w http.ResponseWriter, r *http.Request) {
	soapResponse := `
	<soap:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"><wsse:Security xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" SOAP-ENV:mustUnderstand="1"><ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="SIG-CC9AF4DE3C0731DB8816919945734959406454"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#gost34310-gost34311"/><ds:Reference URI="#id-CC9AF4DE3C0731DB8816919945734959406453"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#gost34311"/><ds:DigestValue>H901gT+G7m5JBH+F8Sq74JCQrw3ZcDYPr7OfXsVn4FY=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>tH6ZyzEZSLo6kp2CGNCVtZLwNDkTyz3Ui+zjtn5m1C/X1wPavDgeO7cHt51gND2+vkQbx4Egq3M5rHn0VHmVGQ==</ds:SignatureValue><ds:KeyInfo Id="KI-CC9AF4DE3C0731DB8816919945734959406451"><wsse:SecurityTokenReference wsu:Id="STR-CC9AF4DE3C0731DB8816919945734959406452"><ds:X509Data><ds:X509IssuerSerial><ds:X509IssuerName>CN=ҰЛТТЫҚ КУӘЛАНДЫРУШЫ ОРТАЛЫҚ (GOST),C=KZ</ds:X509IssuerName><ds:X509SerialNumber>231278386876974093604887536074399274144411123947</ds:X509SerialNumber></ds:X509IssuerSerial></ds:X509Data></wsse:SecurityTokenReference></ds:KeyInfo></ds:Signature></wsse:Security></SOAP-ENV:Header><SOAP-ENV:Body xmlns:ns1="http://bip.bee.kz/SyncChannel/v10/Types" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" wsu:Id="id-CC9AF4DE3C0731DB8816919945734959406453"><ns1:SendMessageResponse xmlns:ns1="http://bip.bee.kz/SyncChannel/v10/Types"><response><responseInfo><messageId>b9e9af8a-6573-4466-b8a7-edb63e3a4a34</messageId><responseDate>2023-08-14T12:29:32.263+06:00</responseDate><status><code>Success</code><message>ÐžÐš</message></status><sessionId>{5a7540fd-3b13-4acc-848d-5bf33c94ef4b}</sessionId></responseInfo><responseData xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><data xmlns="http://www.w3.org/2001/XMLSchema-instance" type="xsd:string">&lt;?xml version="1.0" encoding="UTF-8" standalone="yes"?>&lt;ns2:response xmlns:ns2="http://shep.nitec.kz/">&lt;requestNumber>10fdc2fa-e9ba-4169-907d-71a272bd78c5&lt;/requestNumber>&lt;status>&lt;code>SUCCESS&lt;/code>&lt;messageRu>OK&lt;/messageRu>&lt;messageKz>OK&lt;/messageKz>&lt;/status>&lt;/ns2:response></data></responseData></response></ns1:SendMessageResponse></SOAP-ENV:Body></soap:Envelope>`
	
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprint(w, soapResponse)
}

func main() {
	port := "80"
	path := "/bip-sync-wss-gost/"
	
	http.HandleFunc(path, soapHandler)
	
	fmt.Printf("Starting test SOAP server on port %s...\n", port)
	err := http.ListenAndServe("", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
