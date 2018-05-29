
wget "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-$(go env GOOS)-x86_64-1.8.0.tar.gz" 
sudo tar xvfz libtensorflow-${TF_TYPE}-$(go env GOOS)-x86_64-1.8.0.tar.gz -C /usr/local

