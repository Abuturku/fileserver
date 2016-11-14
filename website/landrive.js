//todo not hard coded!
var depthOfFoldersUnderRoot = 0;
//todo not hard coded!
var username = "Max Muster";

var folderData;
var currentFolderPath = "";

var folderBacklog = [];
var currentFolder;
var folderForwlog = [];

//todo catch server response instead of using hardcoded data
function loadFolderData(data){
	folderData = data;
	username = folderData.Name;
}

function deactivateButton(sId){
	//make sure that the class won't be duplicated
	document.getElementsByClassName(sId)[0].classList.remove("inactive_icon");
	document.getElementsByClassName(sId)[0].classList.add("inactive_icon"); 

}

function activateButton(sId){
	document.getElementsByClassName(sId)[0].classList.remove("inactive_icon");
}

function generateFolderStructure(){
		//show root folder
		var rootHtml = '<div class="folderRoot" onclick="onclickFolderSelected(this, event)"><span id="homeTitle">Home of ';
			rootHtml += username + '</span></div>';
		document.getElementById("folderStructure").innerHTML = rootHtml;
		//show deep structure via recursion
		var foldersHtml = readFolderStructureRec(folderData.Folders, 1);
		document.getElementById("folderStructure").innerHTML += foldersHtml;
	}

	function readFolderStructureRec(childFolders, depth){
		//exit condition 
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
	
window.onload = function () {
	//catch server response
	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
			loadFolderData(JSON.parse(xmlhttp.responseText));
			generateFolderStructure();
			folderSelected(document.getElementsByClassName("folderRoot")[0],null);
		}
	}

	xmlhttp.open("GET", "/getFolderStruct", true);
	xmlhttp.send();
}

/**backup
function searchCurrentFolderObjectRec(childFolders){
	var folderObj;
	for (var i = 0; i < childFolders.length; i++){
		var nameParts = childFolders[i].Name.split("/");
		var nameOfCurrFolder = nameParts.pop();
		if (nameOfCurrFolder === document.getElementById("selectedFolder").children[0].innerHTML){
			return childFolders[i];
		} else {
			var currTest = searchCurrentFolderObjectRec(childFolders[i].Folders);
			if(currTest != undefined){
				folderObj = currTest;
			}
		}
	}
	return folderObj;
}*/

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
/*
function searchCurrentFolderObjectViaPath(){
	
	//todo fix
	var folderpath = currentFolderPath;
	
}*/

function getCurrentFolderObject(){
	if(document.getElementById("selectedFolder").children[0].innerHTML === "Home of "+folderData.Name){
		return folderData;
	}
	return searchCurrentFolderObjectRec(folderData.Folders);
	//return searchCurrentFolderObjectViaPath();
}

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

function onclickDownloadFile(form){
	form.submit();
}

function onclickDeleteFile(form){
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
	
function onFileSelectedForUpload(form){
	form.submit();
}

function onclickDeleteFolder(form){
	var b = confirm("Are you sure that you want to delete the folder with all it's subfolders and files?");
	if (b == true) {
		form.submit();
	} else {
		alert("Deletion of folder cancelled.");
	}
}

function cancelPwChange(){
	var cpw_fields = document.getElementsByClassName("cpw_input");
	for (var i = 0; i<cpw_fields; i++){
		cpw_fields[i].value = "";
	}
	document.getElementById("pwChngDialog").classList.add("hidden");
}

function onclickChangePw(){
	document.getElementById("pwChngDialog").classList.remove("hidden");
}

function onclickNewFolder(form){
	var newFolderName = prompt("Name of the new Folder:", "Example Folder");
	document.getElementById("newFolderName").value = newFolderName;
	form.submit();
}