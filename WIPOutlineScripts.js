var dom = {};
dom.sel = function(selector) {
    return document.querySelectorAll(selector);
}
dom.create = function(element) {
    return document.createElement(element);
}
// For all nodes, we need to two things for now
// WIP
// - An edit/submit changes flow, clicking/tapping 
// on a node should start editing it, and while the node is in
// edit mode, there should be some kind of "save" button for it
//
// TODO: Each node should also have an expand/collapse button
//
// - And then start saving changes to the backend
// TODO: Add UI for creating new nodes
// TODO: Add UI for deleting nodes
// TODO: Add UI for re-ordering and re-paretning nodes

var werewolf = (function() {
    var exports = {};
    exports.nodeSaveClick = function(id) {
        return function(event) {
            // Kick off the reqwest
            // Set 
            var button = event.target;
            var outer = event.target.parentNode;
            var input = nodeInputById(id);
            // TODO: Clean this up on the server side
            var content = input.value.trim();
            button.innerHTML = "Saving...";
            reqwest({
                url: "/node/"+id+"/edit",
                method: "post",
                data: { content: content },
            }).then(function() {
                button.remove();
                input.remove();
                outer.innerHTML = content;
            }).fail(function() {
                outer.innerHTML = "Unable to save node: " + content;
            });
        };
    };
    function nodeInputById(nodeNum) {
        return dom.sel(inputId(nodeNum).sel)[0];
    }
    function nodeButtonById(nodeNum) {
        return dom.sel(saveId(nodeNum).sel)[0];
    }
    function inputId(nodeNum) {
        var sel = "node-"+nodeNum+"-field";
        return { attr: 'id="'+sel+'"', sel: "#"+sel};
    }
    function saveId(nodeNum) {
        var sel = "node-"+nodeNum+"-save";
        return { attr: 'id="'+sel+'"', sel: "#"+sel};
    }

    exports.nodeEditClick = function(event) {
        var child = event.target;
        var id = event.target.dataset.id;
        // Keep this from firing twice
        event.target.removeEventListener('click', exports.nodeEditClick);
        var content = child.innerHTML.trim();
        child.innerHTML = 
            '<input '+inputId(id).attr+' type="text" value="'+content+'"><button '+saveId(id).attr+'>Save</button>' ;
        // Try to fire this off after the DOM has been updated?
        //  setTimeout(function() {
        nodeButtonById(id).addEventListener('click', exports.nodeSaveClick(id));
        // }, 1);
    }
    return exports;
})();


window.onload = function() {
    var nodes = dom.sel(".outline-node .outline-node-content");
    for( var i=0; i< nodes.length;i++) {
        var n = nodes[i];
        n.addEventListener('click', werewolf.nodeEditClick);
    }
}

