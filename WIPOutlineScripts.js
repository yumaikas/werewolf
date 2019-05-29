var dom = {};
dom.sel = function(selector) {
    return document.querySelectorAll(selector);
}

dom.id = function(id) {
    return document.getElementById(id);
}

dom.foreach = function (sel, fn) {
    var nodes = dom.sel(sel);
    for(var i=0; i< nodes.length;i++) {
        var n = nodes[i];
        fn(n, i);
    }
}
dom.create = function(element) {
    return document.createElement(element);
}

// A simple DOM-templating tool
var t = {};

var atr = function() {

    var self = {};
    var pairs = [];
    var ATTRIBUTE_HAS_VALUE = 0;
    var ATTRIBUTE_IS_VOID = 1;

    function attributesWithValues(...names) {
        for (var i = 0; i < names.length;i++) {
            let name = names[i];
            self[names[i]] = function (val) {
                pairs.push({name: name, value: val, type: ATTRIBUTE_HAS_VALUE});
                return self;
            };
        }
    };
    attributesWithValues("id", "type", "value", "class", "for", 
        "name", "cols", "rows", "formmethod", "formaction",
        "form", "action", "accept", "href", "enctype", "rel", "src");
    function voidAttributes (...names) {
        for (var i = 0; i < names.length;i++) {
            let name = names[i];
            self[name[i]] = function () {
                pairs.push({name: name, type: ATTRIBUTE_IS_VOID});
                return self;
            };
        }
    };
    voidAttributes("multiple");
    self.render = function() {
        var buf = "";
        for (var i=0;i<pairs.length;i++) {
            if (pairs[i].type === ATTRIBUTE_IS_VOID) {
                buf += pairs[i].name + " ";
            } else if (pairs[i].type === ATTRIBUTE_HAS_VALUE) {
                buf += pairs[i].name + "=";
                buf += '"' + pairs[i].value + '" ';
            } else {
                console.error("Attribute type isn't set!");
            }
        }
        return buf;
    }

    return self;
};

(function() {
    function voidTag(name, atr) {
        var buf = "<" + name + " ";
        buf += atr.render();
        buf += "/>";
        return buf;
    }

    function nestableTag(name, atr, inner) {
        var buf = "<" + name + " ";
        buf += atr.render();
        buf += ">";
        for (var i=0; i<inner.length; i++) {
            if (typeof inner[i] === "object") {
                buf += inner[i].html();
            } else if (typeof inner[i] === "string") {
                buf += inner[i];
            } else {
                console.error("Unknown type of nested tag!")
            }
        }
        buf +="</"+name+">";
        return buf;
    }

    function addNestedTags(...names) {
        for (var i=0; i<names.length;i++) {
            let name = names[i];
            t[names[i]] = function(atr, ...inner) {
                return {html: function(){
                    return nestableTag(name, atr, inner);
                }};
            }
        }
    }

    function addVoidTags(...names) {
        for (var i=0; i<names.length;i++) {
            let name = names[i];
            t[names[i]] = function(atr, ...inner) {
                return {html: function(){
                    return voidTag(name, atr, inner);
                }};
            }
        }
    }

    addVoidTags("input", "br", "hr", "link");
    addNestedTags("div", "button", "details", "summary", "span", "style",
       "h1", "h2", "h3", "h4", "h5", "title", "a", "p", "form", "label", "script");

})();

// For all nodes, we need to two things for now
// on a node should start editing it, and while the node is in
// edit mode, there should be some kind of "save" button for it
//
// TODO: Add UI for creating new nodes
// - WIP: There's a start on the back end for this
// TODO: Add UI for deleting nodes
// TODO: Add UI for re-ordering and re-paretning nodes
// TODO: Figure out how to add a global shortcut (maybe Alt+/?) for entering command mode
// TODO: Figure out how to make a navigation mode.

var lunacy = (function() {
    var exports = {};
    exports.createNodeFrom = function(parentId, outline_order, content) {
        return request({
            url: "",
            method: "post",
            data: {
                content: content,
                parent_id: parentId,
                outliner_order, outline_order,
            }
        });
    }
    exports.saveNodeChanges = function(id, content) {
        return reqwest({
            url: "/node/"+id+"/edit",
            method: "post",
            data: {content: content},
        });
    };
    return exports;
})()

// The part you're interacting with
var werewolf = (function() {
    var exports = {};

    // Nodes in limbo have a custom state that means they shouldn't be getting general click
    // events, since they'll have their own set of handlers wired up.
    var nodesInLimbo = new Set();

    function nodeInputById(nodeNum) {
        return dom.sel(inputId(nodeNum).sel)[0];
    }

    function saveId(nodeNum) {
        var sel = "node-"+nodeNum+"-save";
        return { id: sel, sel: "#"+sel};
    }
    function nodeButtonById(nodeNum) {
        return dom.sel(saveId(nodeNum).sel)[0];
    }
    function inputId(nodeNum) {
        var sel = "node-"+nodeNum+"-field";
        return { id: sel, sel: "#"+sel};
    }

    exports.expandAll = function() {
        var elems = dom.sel("details");
        for(var i =0; i < elems.length; i++) {
            elems[i].open = true;
        }
    }
    exports.collapseAll = function() {
        var elems = dom.sel("details");
        for(var i =0; i < elems.length; i++) {
            elems[i].open = false;
        }
    }
    var nodeClickFunc = function(event) {
        console.error("unset nodeClickFunc!");
    };

    exports.nodeClick = function(event) {
        // Keep handlers from having to do this over and over.
        event.preventDefault();
        // Make sure that we're not firing for every subnode
        if (!nodesInLimbo.has(event.target)&& event.target.matches(".outline-node-content")) {
            nodeClickFunc.apply(this, arguments);
        }
    }

    exports.flashMessage = function(message) {
        var messageBar = dom.id("message-bar");
        messageBar.classList.remove("hidden");
        messageBar.innerHTML = message;
        dom.id("command-bar").classList.add("hidden");
    }
    exports.clearFlash = function() {
        var messageBar = dom.id("message-bar");
        messageBar.classList.add("hidden");
        messageBar.innerHTML = "";
        dom.id("command-bar").classList.remove("hidden");
    }

    // An attempt an setting tap "modes"
    exports.setTapFunc = function(handler) {
        nodeClickFunc = handler;
    }
    exports.beginCreateNode = function() {
        exports.setTapFunc(exports.createNodeClick);
        exports.flashMessage("Where do you want to create the node?");
    };

    var newNodeSeq = 0;
    exports.createNodeClick = function(event) {
        exports.clearFlash();
        var child = event.target;
        var id = event.target.dataset.id;
        // Find the number of direct children of the node in question.
        var content = child.innerHTML.trim();
        newNodeSeq++;
        var newNodeNum = newNodeSeq;
        var inputId = "new-node-" + newNodeSeq;
        var createBtnId = "save-new-node-" + newNodeSeq;
        var portDivId  = "new-node-div-" + newNodeSeq;

        // So, we need to 
        child.innerHTML = content + 
            t.div(atr().id(portDivId),
                    t.label(atr(), "New node:", 
                        t.input(atr().type("text").id(inputId)),
                        t.button(atr().id(createBtnId), "Create!"
                            ))).html();

        var input = dom.id(inputId);
        var saveBtn = dom.id(createBtnId);
        var portDiv = dom.id(portDivId);
        
        var clickHandler = exports.requestNodeCreation(portDiv, input, saveBtn, id);

        input.addEventListener('keyup', eatSpace);
        saveBtn.addEventListener('click', clickHandler);
    };

    exports.requestNodeCreation = function(portDiv, input, button, parentId) {
        var handler = function(event) {
            // Keep this from bubbling up to the summary/details element
            event.preventDefault();
            var numChildren = dom.sel("#node-"+ parentId +">details").length;
            reqwest({
                url:"/node/create/",
                method: "post",
                data: {
                    content: input.value,
                    parent_id: parentId,
                    outline_order: numChildren + 1
                }
            }).then(function(data) {
                // For now, reload the page.
                location.reload();
                // TODO
            }).fail(function(err) {
                console.error(err);
            })
            button.removeEventListener('click', handler);
            input.removeEventListener('keyup', eatSpace);
            portDiv.delete();
        }
        return handler;
    }

    function eatSpace(event) {
        if (event.code === "Space") {
            event.preventDefault();
        }
    }

    exports.nodeEditClick = function(event) {
        var child = event.target;
        var id = event.target.dataset.id;
        // Clean this up?
        var content = child.innerHTML.trim();
        child.innerHTML = t.input(
                atr().id(inputId(id).id).
                type("text").value(content)
                ).html() +
            t.button(atr().id(saveId(id).id), "Save").html();

        //'<input '+inputId(id).attr+' type="text" value="'+content+'"><button '+saveId(id).attr+'>Save</button>' ;
        // Try to fire this off after the DOM has been updated?
        nodeButtonById(id).addEventListener('click', exports.nodeSaveClick(child, id));
        nodeInputById(id).addEventListener('keyup', eatSpace);
        nodesInLimbo.add(child);
    }

    exports.nodeSaveClick = function(inLimbo, id) {
        return function(event) {
            // Kick off the reqwest
            event.preventDefault();
            var button = event.target;
            var outer = event.target.parentNode;
            var input = nodeInputById(id);
            // TODO: Clean this up on the server side
            var content = input.value.trim();
            button.innerHTML = "Saving...";
            lunacy.saveNodeChanges(id, content).then(function() {
                button.remove();
                input.remove();
                outer.innerHTML = content;
                nodesInLimbo.delete(inLimbo);
            }).fail(function() {
                outer.innerHTML = "Unable to save node: " + content;
                nodesInLimbo.delete(inLimbo);
            });
        };
    };
    nodeClickFunc = exports.nodeEditClick;
    return exports;
})();

window.onload = function() {
    dom.foreach(".outline-node .outline-node-content", function(n) {
        n.addEventListener('click', werewolf.nodeClick);
    });
}
