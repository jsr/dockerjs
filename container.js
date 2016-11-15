function Container( image, name ) {
	this.image = image;
	this.name  = name; 
	result = create({'Image':this.image}, null, null, this.name);
	this.ID = result.ID
}

Container.prototype.run = function() {
	result = run(this.ID);
}

Container.prototype.on = function(eventtype, callback) {
	listen( this.ID, eventtype, callback );
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
		return container.ID.indexOf(id) == 0;
	})[0]

	if( c == null ) {
		return null 
	}

	result = new Container( c.Image, c.Names[0] );
	result.ID = c.ID;
	return result
}

