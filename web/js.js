//Materialize initialisieren

init()
var initiatedConns = false;
var initiatedSecrets = false;
function init(){
	jQuery.getJSON("/api/secrets", {}, buildSecretTable);
	jQuery.getJSON("/api/connections", {}, buildConnectionTable);
}
function initiateCollapsible(){
	if (initiatedConns && initiatedSecrets) {
		M.AutoInit()
	}
}
function buildConnectionTable(data){
	var elem = jQuery("#connectionscollapsible")
	var html = "";
	jQuery.each(data, function( index, v){
		html += "<li>";
		html += "<div class=\"collapsible-header\">" + v.Path + "</div>";
		html += "<div class=\"collapsible-body\">";
		html += "<div class=\"row\">";
		html += "<div class =\"input-field col s4 center\">";
		html += "<input value=\"" + v.Path + "\" placeholder=\"" + v.Path + "\" id=\"connection-path-" + v.Path + "\" type=\"text\" disabled>";
		html += "<label class=\"active\" for=\"connection-path-" + v.Path + "\">Pfad zur Datei</label>";
		html += "</div>";
		html += "<div class =\"input-field col s4 center\">";
		html += "<input value=\"" + v.LocalAddrs + "\" placeholder=\"" + v.LocalAddrs + "\" id=\"connection-local-ip-" + v.Path + "\" type=\"text\">";
		html += "<label class=\"active\" for=\"connection-local-ip-" + v.Path + "\">Lokale Public IP</label>";
		html += "</div>";
		html += "<div class =\"input-field col s4 center\">";
		html += "<input value=\"" + v.RemoteAddrs + "\" placeholder=\"" + v.RemoteAddrs + "\" id=\"connection-remote-ip-" + v.Path + "\" type=\"text\">";
		html += "<label class=\"active\" for=\"connection-remote-ip-" + v.Path + "\">Remote Public IP</label>";
		html += "</div>";
		html += "<div>";
		html += "<div class=\"row\">";
		html += "<div class =\"input-field col s4 center\">";
		html += "<input value=\"" + v.Version + "\" placeholder=\"" + v.Version + "\" id=\"connection-version-" + v.Path + "\" type=\"text\">";
		html += "<label class=\"active\" for=\"connection-version-" + v.Path + "\">IKE Version</label>";
		html += "</div>";
		html += "<div class =\"input-field col s4 center\">";
		html += "<input value=\"" + v.Proposals + "\" placeholder=\"" + v.Proposals + "\" id=\"connection-proposals-" + v.Path + "\" type=\"text\">";
		html += "<label class=\"active\" for=\"connection-proposals-" + v.Path + "\">Proposal for IKE</label>";
		html += "</div>";
		html += "<div class =\"input-field col s4 center\">";
		html += "<input value=\"" + v.ChildProposals + "\" placeholder=\"" + v.ChildProposals + "\" id=\"connection-child-proposals-" + v.Path + "\" type=\"text\">";
		html += "<label class=\"active\" for=\"connection-child-proposals-" + v.Path + "\">Proposals for Children</label>";
		html += "</div>";
		html += "</div>";
		html += "<div class=\"row\">";
		html += "<div class =\"input-field col s12 center\">";
		html += "<input value=\"" + v.LocalTS + "\" placeholder=\"" + v.LocalTS + "\" id=\"connection-local-ts-" + v.Path + "\" type=\"text\">";
		html += "<label class=\"active\" for=\"connection-local-ts-" + v.Path + "\">TrafficSelector for LocalSubnets</label>";
		html += "</div>";
		html += "</div>";
		html += "<div class=\"row\">";
		html += "<div class =\"input-field col s12 center\">";
		html += "<input value=\"" + v.RemoteTS + "\" placeholder=\"" + v.RemoteTS + "\" id=\"connection-remote-ts-" + v.Path + "\">";
		html += "<label class=\"active\" for=\"connection-remote-ts-" + v.Path + "\">TrafficSelector for RemoteSubnets</label>";
		html += "</div>";
		html += "</div>";
		html += "<div class=\"row\">";
		html += "<div class=\"switch\">";
		var initiator = v.Initiator == "yes" ? "checked" : "";
		html += "<label>Initiator<input type=\"checkbox\" " + initiator + " id=\"initiator-" + v.Path + "\"><span class=\"lever\"></span></label>";
		html += "</div>";
		html += "<div class=\"switch\">";
		var udp = v.UDPEncap == "yes" ? "checked" : "";
		html += "<label>UDP Encapsulation<input type=\"checkbox\" " + udp + " id=\"udp-encap-" + v.Path + "\"><span class=\"lever\"></span></label>";
		html += "</div>";
		html += "<button type=\"button\" class=\"waves-effect waves-light btn\" onclick=\"onConnectionEdit('" + v.Path + "')\">Edit</button>";
		html += "<button type=\"button\" class=\"waves-effect waves-light btn\" onclick=\"onConnectionLoad('" + v.Path + "')\">Load</button>";
		html += "<button type=\"button\" class=\"waves-effect waves-light btn\" onclick=\"onConnectionUnload('" + v.Path + "')\">Unload</button>";
		html += "<button type=\"button\" class=\"waves-effect waves-light btn\" onclick=\"onConnectionDelete('" + v.Path + "')\">Delete</button>";
		html += "</div>";
		html += "</div>";
		html += "</li>";
	})
	html += elem.html();
	elem.html(html);
	initiatedConns = true
	initiateCollapsible()
}
function getConnectionFromInputs(path){
	return {
		Path: jQuery("#connection-path-" + path).val(),
		LocalAddrs: jQuery("#connection-local-ip-" + path).val(),
		RemoteAddrs: jQuery("#connection-remote-ip-" + path).val(),
		Proposals: jQuery("#connection-proposals-" + path).val(),
		ChildProposals: jQuery("#connection-child-proposals-" + path).val(),
		Version: jQuery("#connection-version-" + path).val(),
		RemoteTS: jQuery("#connection-remote-ts-" + path).val(),
		LocalTS: jQuery("#connection-local-ts-" + path).val(),
		UDPEncap: jQuery("#udp-encap-" + path).val()? "yes" : "no",
		Initiator: jQuery("#initiator-" + path).val()? "yes" : "no"
	};
}
function setConnectionToInputs(c){
	jQuery("#connection-path-" + c.Path).val(c.Path);
	jQuery("#connection-local-ip-" + c.Path).val(c.LocalAddrs);
	jQuery("#connection-remote-ip-" + c.Path).val(c.RemoteAddrs);
	jQuery("#connection-proposals-" + c.Path).val(c.Proposals);
	jQuery("#connection-child-proposals-" + c.Path).val(c.ChildProposals);
	jQuery("#connection-version-" + c.Path).val(c.Version);
	jQuery("#connection-remote-ts-" + c.Path).val(c.RemoteTS);
	jQuery("#connection-local-ts-" + c.Path).val(c.LocalTS);
	jQuery("#udp-encap" + c.Path).val(c.UDPEncap == "yes" ? true : false);
	jQuery("#initiator-" + c.Path).val(c.Initiator == "yes" ? true : false);
}
function callBackendConn(method, path, url){
	var data = getConnectionFromInputs(path);
	var requestUrl = "/api/connections/" + data.Path;
	if (url) {
		requestUrl += url;
	}
	jQuery.ajax({
		type: method,
		url: requestUrl,
		data: JSON.stringify( data )
	}).always(callbackConnections);
}
function callbackConnections(connection, callStatus){
	if (callStatus === "success"){
		setConnectionToInputs(connection);
	}else if (callStatus === "error") {
		debugger;
	}
}
function onConnectionEdit(path){
	var method = "PUT";
	if (path === "new") {
		method = "POST";
	}
	callBackendConn(method, path, undefined);
}
function onConnectionDelete(path){
	callBackendConn("DELETE", path, undefined);
}
function onConnectionLoad(path){
	callBackendConn("PUT", path, "/load");
}
function onConnectionUnload(path){
	callBackendConn("PUT", path, "/unload");
}
function buildSecretTable(data){
	var elem = jQuery("#secretscollapsible")
	var html = "";
	jQuery.each(data, function(index, value){
		html += "<li>";
		//Kopfzeile
		html += "<div class=\"collapsible-header\">";
		html += value.Path;
		html += "</div>";
		
		//Body auf
		html += "<div class=\"collapsible-body\">";
		//Owner
		html += "<div class=\"row\">";
		html += "<div class=\"input-field col s4\">";
		html += "<input value=\"" + value.Owners + "\" placeholder=\"" + value.Owners + "\" id=\"secret-owner-" + value.Path + "\" type=\"text\">";
		html += "<label class=\"active\" for=\"secret-owner-" + value.Path + "\">Owner (Others Public IP)</label>";
		html += "</div>";
		//PSK
		html += "<div class=\"input-field col s4\">";
		html += "<input value=\"" + value.Data + "\" placeholder=\"" + value.Data + "\" id=\"secret-psk-" + value.Path + "\" type=\"text\">";
		html += "<label class=\"active\" for=\"secret-psk-" + value.Path + "\">Preshared Key</label>";
		html += "</div>";
		//Path
		html += "<div class=\"input-field col s4\">";
		html += "<input value=\"" + value.Path + "\" placeholder=\"" + value.Path + "\" id=\"secret-path-" + value.Path + "\" type=\"text\" disabled>";
		html += "<label class=\"active\" for=\"secret-psk-" + value.Path + "\">Pfad zur Datei</label>";
		html += "</div>";
		html += "</div>";
		
		html += "<div class=\"row\">";
		//Buttons
		html += "<button type=\"button\" class=\"waves-effect waves-light btn\" onclick=\"onSecretEdit('" + value.Path + "')\">Edit</button>";
		html += "<button type=\"button\" class=\"waves-effect waves-light btn\" onclick=\"onSecretLoad('" + value.Path + "')\">Load</button>";
		html += "<button type=\"button\" class=\"waves-effect waves-light btn\" onclick=\"onSecretUnload('" + value.Path + "')\">Unload</button>";
		html += "<button type=\"button\" class=\"waves-effect waves-light btn\" onclick=\"onSecretDelete('" + value.Path + "')\">Delete</button>";
		html += "</div>";
		//Body zu
		html += "</div>";
		
		html += "</li>"
		
	})
	html += elem.html();;
	elem.html(html);
	initiatedSecrets = true
	initiateCollapsible()
}
function getSecretFromInputs(path){
	return {
		Path: jQuery("#secret-path-" + path).val(),
		Owners: jQuery("#secret-owner-" + path).val(),
		Data: jQuery("#secret-psk-" + path).val(),
		Typ: "PSK"
	};
}
function setSecretToInputs(secret){
	jQuery("#secret-owner-" + secret.Path).val(secret.Owners)
	jQuery("#secret-psk-" + secret.Path).val(secret.Data)
	jQuery("#secret-path-" + secret.Path).val(secret.Path)
}
function onSecretEdit(path){
	var method = "PUT";
	if (path === "new"){
		method = "POST";
	}
	callBackendSecret(method, path, undefined)
}
function callBackendSecret(method, path, url){
	var data = getSecretFromInputs(path);
	var requestUrl = "/api/secrets/" + data.Path;
	if (url) {
		requestUrl += url
	}
	jQuery.ajax({
		type: method,
		url: requestUrl,
		data: JSON.stringify( data )
	}).always(callbackSecrets)
}
function onSecretDelete(path){
	callBackendSecret("DELETE", path, undefined)
}
function onSecretLoad(path){
	callBackendSecret("PUT", path, "/load");
}
function onSecretUnload(path){
	callBackendSecret("PUT", path, "/unload");
}

function callbackSecrets(secret, callStatus){
	if (callStatus === "success") {
		setSecretToInputs(secret);
	}else if (callStatus === "error") {
		debugger;
	}
}
