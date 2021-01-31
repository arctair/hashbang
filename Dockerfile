FROM openjdk:11-jre
COPY hashbang-*.jar /app/hashbang.jar
EXPOSE 8080
CMD ["java", "-jar", "/app/hashbang.jar"]
