FROM nginx:stable-alpine
RUN rm /usr/share/nginx/html/*.*
COPY index.html /usr/share/nginx/html
COPY . /usr/share/nginx/html/
RUN ls -R /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]