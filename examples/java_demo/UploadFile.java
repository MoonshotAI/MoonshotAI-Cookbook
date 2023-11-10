public FileUploadResult uploadInputStream(InputStream inputStream, String fileName, long cententLength) {
        // 1、封装请求头
        HttpHeaders headers = new HttpHeaders();
        MediaType type = MediaType.parseMediaType("multipart/form-data");
        headers.setContentType(type);
        headers.add("Authorization", KimiAiConstant.BEARER+KimiAiConstant.MOONSHOT_API_KEY);
        // 2、封装请求体
        MultiValueMap<String, Object> param = new LinkedMultiValueMap<>();
        InputStreamResource resource = new InputStreamResource(inputStream){
            @Override
            public long contentLength(){
                return cententLength;
            }
            @Override
            public String getFilename(){
                return fileName;
            }
        };
        param.add("file", resource);
        param.add("purpose", "file-extract");
        // 3、封装整个请求报文
        HttpEntity<MultiValueMap<String, Object>> formEntity = new HttpEntity<>(param, headers);
        // 4、发送请求
        ResponseEntity<String> data = restTemplate.postForEntity(KimiAiConstant.UPLOAD_FILES_URL, formEntity, String.class);
        // 5、请求结果处理
        FileUploadResult fileUploadResult = JSONObject.parseObject(data.getBody(), FileUploadResult.class);
        // 6、返回结果
        return fileUploadResult;
    }