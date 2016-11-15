function Container( image, name ) {
	this.image = image;
	this.name  = name; 
}

Container.prototype.run = function() {
	result = run({Image: image}, null, null, name);
	if( result != null ) {
		this.ID = result.ID
	}
}

Container.all = function() {
	return list()
}

Container.match = function(predicate) {
	result = []
	containers = Container.all();
	for( i = 0; i < containers.length; i++) { 
		if( predicate(containers[i]) ) {
			result.push( containers[i] )
		}
	}
	return result;
}

Container.findByID = function(id) {
	c = Container.match(function(container) {
		return container.ID == id;
	})[0]

	result = new Container( c.Image, c.Names[0] );
	result.ID = c.ID;
	return result
}
