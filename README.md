# urlShortner
The purpose of this project is to learn to work with Docker and Kubernetes. In fact, we are going to containerize a relatively simple project using Docker and deploy it on Kubernetes(minikube) by writing Kubernetes deployment files. To do this project, we need to install docker and minikube on our system. For this purpose, we can get help from the following links:

https://docs.docker.com/get-docker/

https://minikube.sigs.k8s.io/docs/start/

In this project, each user can enter the desired link and the server will deliver a small and summarized link so that the primary link address can be accessed using it. For this purpose, we develop a server with Go language, which is connected to a Redis database that stores unique shortened addresses in it. Since this service is provided for free, the saved addresses have an expiration date and are not permanently accessible.

Each user can enter the desired link and the server will deliver a small and summarized link so that the primary link address can be accessed using it. For this purpose, we develop a server with Go language, which is connected to a redis database that stores unique shortened addresses in it. Since this service is provided for free, the saved addresses have an expiration date and are not permenant.

Our server has two endpoints, one of which is responsible for creating the shortened address and the other is responsible for transferring to the shortened address. For example, by sending a bs HTTP request of the Post method to the following endpoint, the server will send the shortened address to the user.

![image](https://github.com/MahdiTheGreat/urlShortner/assets/47212121/ac499605-93b3-4efa-8aee-da6c93a7eae5)

Suppose the server's answer, in this case, is the shortened address "shorturLat/nsBL5" and by sending a request to the second endpoint, which has the same address shortened, The long initial address is loaded for us, like below:

![image](https://github.com/MahdiTheGreat/urlShortner/assets/47212121/bf933b28-aa6f-49d4-982d-accd4376f680)

Keep in mind that this shortened address is accessible for a limited period of time and then it expires.

# Implementation

In order to not have hard-coded configuration, we store the configurations in a separate config file. These configs include:

- server port
- expiration time
- Database server address
- Database name and password

These configs are used when the main server is going up.

We also use a secret file, which is responsible for storing the name and password of the database. Since this information, especially the password, is sensitive information, it should be stored in Secret.

For implementation, we use a Dockerfile by which the project can be containerized. Finally, by building the Dockerfile, we generate the project image and place it on Dockerhub. To create an image, we use the multistage build technique, in which the first step is to build your project and create an executable file so that this file can be executed in an Alpine container in the second step. The console log of this process can be seen below:

<code>C:\urlShortner>docker build -t mahdithegreat/redis-app:5.0 .[+] Building 67.0s (18/18) FINISHED
 => [internal] load build definition from Dockerfile                                                               0.1s
 => => transferring dockerfile: 32B                                                                                0.1s
 => [internal] load .dockerignore                                                                                  0.1s
 => => transferring context: 2B                                                                                    0.0s
 => [internal] load metadata for docker.io/library/alpine:latest                                                  17.5s
 => [internal] load metadata for docker.io/library/golang:latest                                                  17.5s
 => [auth] library/golang:pull token for registry-1.docker.io                                                      0.0s
 => [auth] library/alpine:pull token for registry-1.docker.io                                                      0.0s
 => [builder 1/6] FROM docker.io/library/golang:latest@sha256:0fa6504d3f1613f554c42131b8bf2dd1b2346fb69c2fc24a312  0.0s
 => [stage-1 1/4] FROM docker.io/library/alpine:latest@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc01182  0.0s
 => [internal] load build context                                                                                  2.1s
 => => transferring context: 10.10kB                                                                               1.9s
 => CACHED [builder 2/6] WORKDIR /app                                                                              0.0s
 => CACHED [builder 3/6] COPY go.mod go.sum ./                                                                     0.0s
 => CACHED [builder 4/6] RUN go mod download                                                                       0.0s
 => [builder 5/6] COPY . .                                                                                         0.4s
 => [builder 6/6] RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o httpServer .                     46.0s
 => CACHED [stage-1 2/4] RUN apk --no-cache add ca-certificates                                                    0.0s
 => CACHED [stage-1 3/4] WORKDIR /root/                                                                            0.0s
 => [stage-1 4/4] COPY --from=builder /app/httpServer .                                                            0.3s
 => exporting to image                                                                                             0.2s
 => => exporting layers                                                                                            0.2s
 => => writing image sha256:69cd7fca87cdcd213ead2d4b2295d306f4a081f735585f168635d9005fe86bee                       0.0s
 => => naming to docker.io/mahdithegreat/redis-app:5.0                                                             0.0s
Use 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them
C:\urlShortner>docker push mahdithegreat/redis-app:5.0
The push refers to repository [docker.io/mahdithegreat/redis-app]
419d089b0938: Pushed                                                                                                    
5f70bf18a086: Layer already exists                                                                                      
789c9843f753: Layer already exists                                                                                      
8d3ac3489996: Layer already exists                                                                                      
5.0: digest: sha256:50b2b0a4b1e5d6f3bc437589d3c2d3f17620f20af153c94981f154435add032f size: 1155</code>

In order to keep the database information persistent in the event of a problem with the respective pods, it is necessary to define a persistent volume for it. As a result, the next step is to create a Persistent Volume and then create a Persistent Volume Claim to use it.
Then we need to write a deployment file that is responsible for preparing and maintaining the database (which uses the password defined in Secret).

Note that for all the deployments in this project, we set the number of replicas to 2. In all of the deployments we reference the config file (in the shape of configMap) and the secret file(in the shape of secret) which are needed to start the servers.

We can also statefulSets rather than deployments since deployments are used to implement stateless apps that do not care which network is using them(do not care about the identity of the user) and do not need permanent storage. Although we solve the need for storage by creating a persistent volume, the problem that remains is that in deployments Pods are interchangeable and don't have unique ids, which is a problem if we want to reach specific pods or directly communicate with the master(the master has the most up to date information). In other words, StatefulSets are more suited to implement stateful apps(that typically include a database, such as this project), however, we can still do this project with deployment and persistent volume. We also must keep in mind that the Persistent Volume Claim created in the previous step is mounted on this deployment. 

The question that may arise is how many pods should be suitable for deployment. In answer, we can say that the appropriate number of pods is relative to the size of the persistent volume and persistent volume claim. Basically, the number of pods multiplied by the size of the persistent volume claim should not outsize the size of the persistent volume.

The last thing we need to create is a service that can be used to access the project and the server that we have developed after creating the said files. Basically, to access the database, a service is needed, which can be used to connect the project and the database.
 
After using the deployment and service files by the kubectl apply command on the minikube cluster, To confirm the creation of pods, services, and deployments and the usage of config and secret files, we use several commands which can be seen in the command log below:

<code>C:\cloud Computing\finalProject>kubectl get pods
NAME                            READY   STATUS    RESTARTS   AGE
alpine-test-7fccc6698f-lnb5f    1/1     Running   0          4h21m
redis-app-6b89cb5d54-cbdzd      1/1     Running   0          4m56s
redis-app-6b89cb5d54-gj2dz      1/1     Running   0          4h21m
redis-master-588d9c5554-7f8vr   1/1     Running   0          4h10m
C:\urlShortner>kubectl get services
NAME                TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
kubernetes          ClusterIP   10.96.0.1       <none>        443/TCP          3d10h
redis-app-service   NodePort    10.106.94.33    <none>        9000:32581/TCP   27h
redis-master        ClusterIP   10.108.74.250   <none>        6379/TCP         2d10h
C:\urlShortner>kubectl get deployments
NAME           READY   UP-TO-DATE   AVAILABLE   AGE
alpine-test    1/1     1            1           4h22m
redis-app      2/2     2            2           4h23m
redis-master   1/1     1            1           4h44m
C:\urlShortner>kubectl get configMaps
NAME                          DATA   AGE
host-port-dbhost-dbport-exp   5      9h
kube-root-ca.crt              1      3d10h
C:\urlShortner>kubectl get secrets
NAME                  TYPE                                  DATA   AGE
default-token-nxsc6   kubernetes.io/service-account-token   3      3d10h
pass-secret           Opaque                                1      28h</code>

keep in mind that for the deployment related to the database we used one pod ,however, we could have used more pods since Redis masters update each other and perhaps in the real world it's better to have several pods for database implementation.

Next, we create two HPA components in the Kubernetes cluster in order to perform auto-scaling operations for both the databases and the pods that are handling requests. The parameters used for autoscaling include average CPU utilization and memory and the reason we chose these parameters is that they are the most basic parameters and both the database and pods handling the requests require both CPU ( in databases for finding the value associated with the key and in pods for handling requests) and memory(in databases for loading the data on to the ram for faster response and in pods for storing the request messages).

<code>C:\urlShortner>kubectl get hpa
NAME                   REFERENCE                 TARGETS                          MINPODS   MAXPODS   REPLICAS   AGE
redis-app-autoscaler   Deployment/redis-app      <unknown>/100Mi, <unknown>/50%   2         5         2          8m19s
redis-db-autoscaler    Deployment/redis-master   <unknown>/100Mi, <unknown>/50%   1         5         1          8m29s</code>

# Testing
We can test our project by using port forwarding of the service created for the project or by using an image with curl capability. For more realism we will use the second method. In the first step, we need to create a docker image based on ubuntu which has the curl command. Then upload the created image on Dockerhub and finally, in order to test it, with the help of the docker rum command, get the image from your Dockerhub and create a container from it, and by sending a curl request to Google.com we test it, which can be see in the command log below:

1)
<code>C:\Users\Mahdi>docker commit 50fb08c1a73c ubuntu-upgraded:1.0
sha256:0ea18b5b83c2fcc5d6ff0fae392899555f29e2d5d497266b61276a284e6e407cC:\Windows\system32>docker push mahdithegreat/ubuntu-upgraded:1.0
C:\cloud Computing\dockerProject\playGround>docker image ls
REPOSITORY               TAG       IMAGE ID       CREATED          SIZE
ubuntu-upgraded          1.0       0ea18b5b83c2   10 minutes ago   207MB
alpine                   latest    0a97eee8041e   10 days ago      5.61MB
docker/getting-started   latest    eb9194091564   11 days ago      28.5MB
ubuntu                   latest    ba6acccedd29   5 weeks ago      72.8MB          
The push refers to repository [docker.io/mahdithegreat/ubuntu-upgraded]
0407c7a29331: Layer already exists                                                                                      
9f54eef41275: Layer already exists                                                                                      
1.0: digest: sha256:279a6bc53ac2de4ba92c1c44a025581b0a6386a2e6553e881080a1ddaf1d3f02 size: 741</code>

2)
<code>C:\cloud Computing\dockerProject\playGround>docker pull mahdithegreat/ubuntu-upgraded:1.0
1.0: Pulling from mahdithegreat/ubuntu-upgraded
Digest: sha256:279a6bc53ac2de4ba92c1c44a025581b0a6386a2e6553e881080a1ddaf1d3f02
Status: Downloaded newer image for mahdithegreat/ubuntu-upgraded:1.0
docker.io/mahdithegreat/ubuntu-upgraded:1.0</code>

3)
<code>root@5cdb8ea19a6b:/# apt-get install curl -y
Reading package lists... Done
Building dependency tree
Reading state information... Done
curl is already the newest version (7.68.0-1ubuntu2.7).
0 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
root@5cdb8ea19a6b:/# curl google.com</code>




