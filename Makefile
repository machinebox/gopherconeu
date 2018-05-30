
dataset: news20.tar.gz		

news20.tar.gz:
	wget http://www.cs.cmu.edu/afs/cs.cmu.edu/project/theo-20/www/data/news20.tar.gz
	tar xvfz news20.tar.gz

wordemb: glove.6B.zip
	unzip glove.6B.zip -d glove.6B
	
glove.6B.zip:
	wget http://nlp.stanford.edu/data/glove.6B.zip

download: dataset wordemb 

install:
	pip2 install -r keras/requirements.txt
