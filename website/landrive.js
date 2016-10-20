//todo not hard coded!
var depthOfFoldersUnderRoot = 3;
//todo not hard coded!
var username = "Max Muster";

var folderBacklog = [];
var currentFolder;
var folderForwlog = [];

function loadFolderData(){
	/*
	var mydata = JSON.parse(data);
	alert(mydata[0].name);
	alert(mydata[0].age);
	alert(mydata[1].name);
	alert(mydata[1].age);
	*/
}

function deactivateButton(sId){
	//make sure that the class won't be duplicated
	document.getElementsByClassName(sId)[0].classList.remove("inactive_icon");
	document.getElementsByClassName(sId)[0].classList.add("inactive_icon"); 

}
function activateButton(sId){
	document.getElementsByClassName(sId)[0].classList.remove("inactive_icon");
}

function loadFilesOfFolder(foldername){
	//todo use real data!
	loadFilesOfFolderWithDummyData(foldername);
}

function loadFilesOfFolderWithDummyData(foldername){
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
	
	