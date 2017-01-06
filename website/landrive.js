/**
  * @author Andreas Schick (2792119), Linda Latreider (7743782), Niklas Nikisch (9364290)
  */

  /**
  * declaration of variables
  */
var depthOfFoldersUnderRoot = 0;
var username = "User";

var folderData;
var currentFolderPath = "";

var folderBacklog = [];
var currentFolder;
var folderForwlog = [];

 /**
 * initialize data of the user, folders and files
 */
 function loadFolderData(data){
	folderData = data;
	username = folderData.Name;
}

 /**
 * changing class of the respective button with the ID "sId"
 * - change accessibility to false
 * - change appearance
 */
function deactivateButton(sId){
	//make sure that the class won't be duplicated
	document.getElementsByClassName(sId)[0].classList.remove("inactive_icon");
	document.getElementsByClassName(sId)[0].classList.add("inactive_icon"); 
}

 /**
 * changing class of the respective button with the ID "sId"
 * - change accessibility to false
 * - change appearance
 */
function activateButton(sId){
	document.getElementsByClassName(sId)[0].classList.remove("inactive_icon");
}

 /**
 * generating html structure of the nested folders beneath the root folder
 */
function generateFolderStructure(){
	//show root folder
	var rootHtml = '<div class="folderRoot" onclick="onclickFolderSelected(this, event)"><span id="homeTitle">Home of ';
		rootHtml += username + '</span></div>';
	document.getElementById("folderStructure").innerHTML = rootHtml;
	//show deep structure via recursion
	var foldersHtml = readFolderStructureRec(folderData.Folders, 1);
	document.getElementById("folderStructure").innerHTML += foldersHtml;
}

 /**
 * uses recursion to transform JSON structure information to html
 * - update depthOfFoldersUnderRoot variable for correct manipulation of the "selected" id
 */
function readFolderStructureRec(childFolders, depth){
	//exit condition: no children left
	if (childFolders.length === 0){
		return "";
	}
	var foldersHtmlTemp = "";
	//recursive calls
	for (var i = 0; i < childFolders.length; i++){
		//generate new folder 
		foldersHtmlTemp += '<div class="folderChild" onclick="onclickFolderSelected(this, event)"><span>';
		var currentName = childFolders[i].Name;
		//only show text after the last '/'
		var nameParts = currentName.split("/");
		foldersHtmlTemp += nameParts.pop() + '</span>';
		
		//test if depth is new maximum
		if (depth > depthOfFoldersUnderRoot){
			depthOfFoldersUnderRoot = depth;
		}
		foldersHtmlTemp += readFolderStructureRec(childFolders[i].Folders, depth+1);
		foldersHtmlTemp += '</div>';
	}
	return foldersHtmlTemp;
}

 /**
 * reading parameters that are added to the url by the server
 */
function getUrlParameter(paramName){
	var result = "-1",
		tmp = [];
	location
		.search.substr(1)
		.split("&")
		.forEach(function (item) {
			tmp = item.split("=");
			if (tmp[0] === paramName){
				result = decodeURIComponent(tmp[1]);
			}
	});
	return result;
}

window.onload = function () {
	 /**
	 * if the change password dialogue failed, the server adds a parameter to the url.
	 * -> the respective problem can be evaluated, so that an error message can be displayed
	 * if there is no issue / if the window is loaded without the former execution of the change password dialogue, nothing happens.
	 */
	//check change pw response
	var message = "none";
	if(getUrlParameter("change")==="pwRepeatFalse"){
		message = "Password change failed.\nThe new passwords did not match.";		
	}
	if(getUrlParameter("change")==="oldPwFalse"){
		message = "Password change failed.\nPlease enter the correct current password.";		
	}
	
	if(message != "none"){
		alert(message);
	}
	
	 /**
	 * when the page loads, the server sends a JSON string with the complete structure that contains information about:
	 * - username of the currently logged in user
	 * - nested folder structure
	 * - files that each folder contains
	 */
	//catch server response
	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
			var jsonString = JSON.parse(xmlhttp.responseText);
			loadFolderData(jsonString);
			generateFolderStructure();
			folderSelected(document.getElementsByClassName("folderRoot")[0],null);
		}
	}

	 /**
	 * define XMLHttpRequest + send
	 */
	xmlhttp.open("GET", "/getFolderStruct", true);
	xmlhttp.send();
}

 /**
 * search the folder and its children
 * - use the currentFolderPath variable to find it
 * - return JSON object of the folder and its children
 */
function searchCurrentFolderObjectRec(childFolders){
	var folderObj;
	for (var i = 0; i < childFolders.length; i++){
		if (childFolders[i].Name === currentFolderPath){
			return childFolders[i];
		} else {
			var currTest = searchCurrentFolderObjectRec(childFolders[i].Folders);
			if(currTest != undefined){
				folderObj = currTest;
			}
		}
	}
	return folderObj;
}

 /**
 * if the current folder is the root folder, return complete structure data
 * else: use recursive search
 */
function getCurrentFolderObject(){
	if(document.getElementById("selectedFolder").children[0].innerHTML === "Home of "+folderData.Name){
		return folderData;
	}
	return searchCurrentFolderObjectRec(folderData.Folders);
}

 /**
 * use the filesize value of the file to create the attribute on the website with the right unit for better readability
 */
function formatFileSize(fileSizeByte){
	var intResult = fileSizeByte;
	var temp;
	nextSize = 1024;
	if (fileSizeByte < nextSize){
		return "" + intResult + " B";
	}
	if (fileSizeByte < nextSize*nextSize){
		temp = "" + fileSizeByte/nextSize;
		intResult = temp.split(".")[0];
		return "" + intResult + " KB";
	}
	nextSize *= 1024;
	if (fileSizeByte < nextSize*nextSize){
		temp = "" + fileSizeByte/nextSize;
		intResult = temp.split(".")[0];
		return "" + intResult + " MB";
	}
	nextSize *= 1024;
	if (fileSizeByte < nextSize*nextSize){
		temp = "" + fileSizeByte/nextSize;
		intResult = temp.split(".")[0];
		return "" + intResult + " GB";
	}
	nextSize *= 1024;
	temp = "" + fileSizeByte/nextSize;
	intResult = temp.split(".")[0];
	return "" + intResult + " TB";
}

 /**
 * generating html structure of the files inside the selected folder
 */
function loadFiles(){
	var fileSpace = document.getElementById("availableFiles");
	var sContent = "";
	var folderObj = getCurrentFolderObject();
	var fileInfo = folderObj.Files;
	
	for(var i=0;i<fileInfo.length;i++){
		//format information
		var fileName = fileInfo[i].Name;
		var fullDate = fileInfo[i].Date.split("T");
		var fileDate = fullDate[0];
		var fileSize = formatFileSize(fileInfo[i].Size);
		
		//create file reference in html
		sContent += '<div class="file" onclick="onclickFileSelected(this)"><span class="fileTitle">';
		sContent += fileName + '</span>';
		sContent += '<div class="fileData"><span class="fileDate">'
		sContent += fileDate + '</span>';
		sContent += '<span class="fileSize">'
		sContent += fileSize + '</span></div></div>'			
	}
	fileSpace.innerHTML = sContent;
}

 /**
 * calculate the path of the selected folder and save it in a global variable
 */
function refreshCurrentFolderPath(){
	var rootName = document.getElementsByClassName("folderRoot")[0].children[0].innerHTML;
	var pathName;
	
	//refresh hidden input fields for form information for server
	var sPath = "";
	var elem = document.getElementById("selectedFolder");
	var folderName = elem.children[0].innerHTML;
	if(folderName !== rootName){
		//search pieces for path
		pathName = [folderName];
		var currParent = elem.parentElement;
		while(currParent !== document.getElementById("folderStructure")){
			pathName.unshift(currParent.children[0].innerHTML);
			currParent = currParent.parentElement;
		}
		//create path with slashes
		var pathLength = pathName.length;
		for (var i = 0; i < pathLength; i++){
			//in front of the first folder no slash!
			if(i !== 0){
				sPath += "/";
			}
			sPath += pathName.shift();
		}
	}
	currentFolderPath = sPath;
}

 /**
 * set the value of the hidden "selected folder" input fields inside the forms of the buttons
 *  in order to pass the information with the http requests
 * uses the global variable currentFolderPath
 * empties the selected file path-fields because no file is selected when a folder gets selected
 */
function refreshHiddenInputFieldsFolders(){
	var inputFields = document.getElementsByClassName("folderPath");
	for (var i = 0; i < inputFields.length; i++){
		inputFields[i].value = currentFolderPath;
	}
	//hidden file path fields must be empty because no file is selected
	var filepathFields = document.getElementsByClassName("filePath");
	for (var j = 0; j < filepathFields.length; j++){
		filepathFields[j].value = "";
	}
}

 /**
 * handles all the necessary actions that need to be taken if a folder gets selected
 * shared functions that need to happen if a folder gets selected, no matter if it is via click or navigation buttons
 */
function folderSelected(elem,event){
	var folderName = elem.children[0].innerHTML;
	
	if (event!== null){
		event.stopPropagation();
	}
	var divs = document.getElementById("folderStructure").children;
	removeFolderIds(divs, depthOfFoldersUnderRoot);
	elem.id = "selectedFolder";
	document.getElementById("folderName").innerHTML = elem.children[0].innerHTML;
	
	refreshCurrentFolderPath();
	refreshHiddenInputFieldsFolders();
	
	//make file buttons unavailable
	deactivateButton("icon_download");
	deactivateButton("icon_delete_file");
	
	//load files of selected folder
	loadFiles();
}

 /**
 * handle all actions that need to be taken if a folder is selected via click
 */
function onclickFolderSelected(elem, event){
	//back navigation
	activateButton("icon_back"); 
	
	//if other folder than root folder is selected, it can be deleted
	if(elem.children[0].id==="homeTitle"){
		deactivateButton("icon_delete_folder");
	} else {
		activateButton("icon_delete_folder");
	}
	
	//save the rootFolder as first element in Backlog
	if(folderBacklog.length === 0){
		folderBacklog.push(document.getElementsByClassName("folderRoot")[0]);
	} else {
		folderBacklog.push(currentFolder);
	}
	
	//save current folder as variable
	currentFolder = elem;
	
	//handle forward log
	deactivateButton("icon_forward");
	folderForwlog = [];
	
	folderSelected(currentFolder,event);
}


 /**
 * remove all IDs of the folders to delete the "selected" attribute
 */
//recursive function to remove marking of past selected folders
function removeFolderIds(divs, remainingFuncCalls){
	if(remainingFuncCalls<=0){
		return;
	}
	
	for (var i=0; i<divs.length; i++){
		//only manipulate divs, not spans!!!
		if(divs[i].tagName === "DIV"){
			//remove ids of current level
			divs[i].removeAttribute("id");
			
			//remove ids of child level
			var divChildren = divs[i].children;
			var remainsNew = remainingFuncCalls - 1;
			if(remainsNew <=0){
				continue;
			}
			removeFolderIds(divChildren, remainsNew);
		}
	}
}

 /**
 * handle click on the navigate back icon
 */
function onclickNavigateBack(){
	//if backlog not empty
	if(folderBacklog.length > 0){
		
		//check if backlog will be empty afterwards
		if(folderBacklog.length === 1){
			deactivateButton("icon_back");
		}
		
		//handle forward log
		activateButton("icon_forward");
		folderForwlog.push(currentFolder);
		
		//save current folder as variable
		currentFolder = folderBacklog.pop();
		
		folderSelected(currentFolder,null);
	}
}

 /**
 * handle click on the navigate forward icon
 */
function onclickNavigateForward(){
	if(folderForwlog.length > 0){
		//check if forwlog will be empty afterwards
		if(folderForwlog.length === 1){
			deactivateButton("icon_forward");
		}
		//handle backward log
		activateButton("icon_back"); 
		folderBacklog.push(currentFolder);
		
		//update current folder 
		currentFolder = folderForwlog.pop();
				
		folderSelected(currentFolder,null);
	}
}

 /**
 * handle all actions that need to be taken if a file is selected
 * set the value of the hidden "selected file" input fields inside the forms of the buttons
 *  in order to pass the information with the http requests
 */
function onclickFileSelected(elem){
	var fileName = elem.children[0].innerHTML;
	
	//unmark all
	var allFiles = document.getElementById("availableFiles").children;
	for (var i = 0; i < allFiles.length; i++){
		allFiles[i].removeAttribute("id");
	}
	//mark selected
	elem.id = "selectedFile";
	
	//make file buttons available
	activateButton("icon_download");
	activateButton("icon_delete_file");
	
	//set filePath
	var pathString = currentFolderPath + "/" + fileName;
	var filepathFields = document.getElementsByClassName("filePath");
	for (var j = 0; j < filepathFields.length; j++){
		filepathFields[j].value = pathString;
	}
}

 /**
 * handle downloading of a file
 */
function onclickDownloadFile(form){
	var buttonClasses = form.children[0].getAttribute("class").split(" ");
	var isInactive = false;
	for(var i = 0; i < buttonClasses.length; i++){
		if(buttonClasses[i] === "inactive_icon"){
			isInactive = true;
		}
	}
	
	if(!isInactive){
		form.submit();
	}
}

 /**
 * handle deletion of a file with a prompt-dialogue
 */
function onclickDeleteFile(form){
	var buttonClasses = form.children[0].getAttribute("class").split(" ");
	var isInactive = false;
	for(var i = 0; i < buttonClasses.length; i++){
		if(buttonClasses[i] === "inactive_icon"){
			isInactive = true;
		}
	}
	
	if(!isInactive){
		var b = confirm("Are you sure that you want to delete the file?");
		if (b == true) {
			//make delete file button unavailable
			deactivateButton("icon_delete_file");
			deactivateButton("icon_download");
			
			form.submit();
		} else {
			alert("Deletion of file cancelled.");
		}
	}
}

 /**
 * upload a file after it has been selected
 */
function onFileSelectedForUpload(form){
	var pathString = form.children[0].children[0].value;//.replace(/\\\\/g, '\\');
	//replace(/\\/g,"/")
//	= form.children[0].children[0].value.replace("\\\\", "\\");
	/*for(var i = 0; i < pathString.length; i ++){
		pathString = pathString.replace("\\", String.fromCharCode(92));
	}*/
	var nameParts = pathString.split("\\");
	var documentName = nameParts.pop();
	document.getElementById("fileNameUpload").value = documentName;
	form.submit();
}

 /**
 * handle deletion of a folder with a prompt-dialogue
 */
function onclickDeleteFolder(form){
	var buttonClasses = form.children[0].getAttribute("class").split(" ");
	var isInactive = false;
	for(var i = 0; i < buttonClasses.length; i++){
		if(buttonClasses[i] === "inactive_icon"){
			isInactive = true;
		}
	}
	
	if(!isInactive){
		var b = confirm("Are you sure that you want to delete the folder with all its subfolders and files?");
		if (b == true) {
			form.submit();
		} else {
			alert("Deletion of folder cancelled.");
		}
	}
}

 /**
 * cancel the change password dialogue and reset the visibility of the change-password-area
 */
function cancelPwChange(){
	var cpw_fields = document.getElementsByClassName("cpw_input");
	for (var i = 0; i<cpw_fields; i++){
		cpw_fields[i].value = "";
	}
	document.getElementById("pwChngDialog").classList.add("hidden");
}

 /**
 * when a http request to the server is sent, make the change-password-area hidden again
 */
function onclickChangePw(){
	document.getElementById("pwChngDialog").classList.remove("hidden");
}

 /**
 * handle error message if a user tries to change his/her password to an empty string
 */
function emptyChangePw(){
	var p1 = document.getElementById("cpw_np1").value;
	var p2 = document.getElementById("cpw_np2").value;
	if(p1 === "" || p2 === ""){
		alert("A password must have at least 1 character.\nPlease fill both \"New password\" fields.");
		return false;
	}
}

 /**
 * handle creation of a new folder with the help of a prompt
 */
function onclickNewFolder(form){
	var newFolderName = prompt("Name of the new Folder:", "Example Folder");
	if(newFolderName!= null){
		var newFolderPath = document.getElementById("newFolderParentPath").value + "/" + newFolderName;
		document.getElementById("newFolderName").value = newFolderPath;
		form.submit();
	}
}