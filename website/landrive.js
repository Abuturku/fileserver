//todo not hard coded!
var depthOfFoldersUnderRoot = 0;
//todo not hard coded!
var username = "Max Muster";

var folderData;
loadFolderData();

var folderBacklog = [];
var currentFolder;
var folderForwlog = [];

//todo catch server response instead of using hardcoded data
function loadFolderData(){
	folderData = {
		"Name": "Andy",
		"Files": [{
			"Name": ".DS_Store",
			"Date": "2016-10-25T12:43:55+02:00",
			"Size": 6148
		}, {
			"Name": "TestDatei.txt",
			"Date": "2016-10-20T15:22:26+02:00",
			"Size": 0
		}],
		"Folders": [{
			"Name": "AndererTest",
			"Files": [{
				"Name": ".DS_Store",
				"Date": "2016-10-25T12:44:06+02:00",
				"Size": 6148
			}],
			"Folders": [{
				"Name": "AndererTest/O1",
				"Files": [],
				"Folders": []
			}, {
				"Name": "AndererTest/O2",
				"Files": [],
				"Folders": []
			}]
		}, {
			"Name": "TestOrdner",
			"Files": [{
				"Name": ".DS_Store",
				"Date": "2016-10-25T12:43:44+02:00",
				"Size": 6148
			}, {
				"Name": "Test",
				"Date": "2016-10-20T16:19:34+02:00",
				"Size": 0
			}, {
				"Name": "TestDateiZwei.txt",
				"Date": "2016-10-20T15:22:36+02:00",
				"Size": 0
			}],
			"Folders": [{
				"Name": "TestOrdner/Test2",
				"Files": [{
					"Name": "Test",
					"Date": "2016-10-20T16:19:34+02:00",
					"Size": 0
				}],
				"Folders": []
			}]
		}]
	}
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

window.onload = function () {
	function generateFolderStructure(){
		//show root folder
		var rootHtml = '<div id="folderRoot" onclick="onclickFolderSelected(this, event)"><span id="homeTitle">Home of ';
			rootHtml += username + '</span></div>';
		document.getElementById("folderStructure").innerHTML = rootHtml;
		//show deep structure via recursion
		var foldersHtml = readFolderStructureRec("", folderData.Folders, 1);
		document.getElementById("folderStructure").innerHTML += foldersHtml;
	}

	function readFolderStructureRec(foldersHtml, childFolders, depth){
		numberOfChildren = childFolders.length;
		//exit condition 
		if (numberOfChildren === 0){
			return "";
		}
		//recursive calls
		for (var i = 0; i < numberOfChildren; i++){
			var nameOfCurrChild = childFolders[i].Name;
			//generate new folder 
			foldersHtml += '<div class="folderChild" onclick="onclickFolderSelected(this, event)"><span>';
			foldersHtml += nameOfCurrChild + '</span>';
			
			//test if depth is new maximum
			if (depth > depthOfFoldersUnderRoot){
				depthOfFoldersUnderRoot = depth;
			}
			foldersHtml += readFolderStructureRec(foldersHtml, childFolders[i].Folders, depth+1);
			foldersHtml += '</div>';
		}
		return foldersHtml;
	}

	generateFolderStructure();
}

function loadFilesDummy(foldername){
	var filesInFolder = document.getElementById("availableFiles");
	var sContent = "";
	var fileInfo = JSON.parse(files);
	
	console.log(fileInfo);
	
	for(var i=0;i<fileInfo.length;i++){
		if(fileInfo[i].fileIn===foldername){
			//create file reference in html
			sContent += '<div class="file" onclick="onclickFileSelected(this)"><span class="fileTitle">';
			sContent += fileInfo[i].fileName + '</span>';
			sContent += '<div class="fileData"><span class="fileDate">'
			sContent += fileInfo[i].fileDate + '</span>';
			sContent += '<span class="fileSize">'
			sContent += fileInfo[i].fileSize + '</span></div></div>'			
		}
	}
	filesInFolder.innerHTML = sContent;
}

//TODO: input-Feld folderPath anpassen bei jedem Ordnerwechsel; alles unterhalb root ||| das selbe für delete mit hidden input
function folderSelected(elem,event){
	var folderName = elem.children[0].innerHTML;
	//load files of selected folder
	/*/todo not hardcoded!
	folderName="Home";
	loadFilesOfFolder(folderName);*/
	
	if (event!== null){
		event.stopPropagation();
	}
	var divs = document.getElementById("folderStructure").children;
	removeFolderIds(divs, depthOfFoldersUnderRoot);
	divs[0].id = "folderRoot";
	elem.id = "selectedFolder";
	document.getElementById("folderName").innerHTML = elem.children[0].innerHTML;
	
	//make file buttons unavailable
	deactivateButton("icon_download");
	deactivateButton("icon_delete_file");
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
		folderBacklog.push(document.getElementById("folderRoot"));
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
				return;
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
	//get name with (elem.children[0].innerHTML);
	//unmark all
	var allFiles = document.getElementById("availableFiles").children;
	for (var i = 0; i<allFiles.length; i++){
		allFiles[i].removeAttribute("id");
	}
	//mark selected
	elem.id = "selectedFile";
	
	//make file buttons available
	activateButton("icon_download");
	activateButton("icon_delete_file");
}

function onclickDownloadFile(){
	alert("TODO: download " + document.getElementById("selectedFile").children[0].innerHTML);
}

function onclickDeleteFile(){
	alert("TODO: delete " + document.getElementById("selectedFile").children[0].innerHTML);
	//make delete file button unavailable
	deactivateButton("icon_delete_file");
	deactivateButton("icon_download");
}
	
function onFileSelectedForUpload(form){
	form.submit();
}