// Create namespace for our app
var app = {};

window.addEvent("domready", function() {
	var boardElement = $("board");

	if (!boardElement) {
		console.error("app: HTML element with id 'board' could not be found.");
		return;
	}

	var board = app.views.board.createInstance(
		boardElement,
		app.config.core.board.columnMaxWidth,
		app.config.core.board.columnMarginLeft
	);
	board.initialize();
	board.rebuild();

	var boardControls = app.views.boardControls.createInstance();
	var boardControlsElement = boardControls.create();
	$(document.body).grab(boardControlsElement);

	var subredditPickerLauncher = new app.ui.SubredditPickerLauncher().toElement();
	boardElement.grab(subredditPickerLauncher, "before");

	window.fireEvent("app.views.boardControls.userDidAskForImages");
});