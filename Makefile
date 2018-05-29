


dataset_test:
	go test -v -run=TestDataset

dataset: news20.tar.gz
	tar xvfz news20.tar.gz	

news20.tar.gz:
	wget http://www.cs.cmu.edu/afs/cs.cmu.edu/project/theo-20/www/data/news20.tar.gz


wordemb: glove.6B.zip
	unzip glove.6B.zip
glove.6B.zip:
	wget http://nlp.stanford.edu/data/glove.6B.zip


install:
	pip2 install ./keras/requirements.txt

.PHONY: dataset_test