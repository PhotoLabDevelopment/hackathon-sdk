package me.photolab;

import com.google.gson.Gson;
// import okhttp3.*;
import okhttp3.OkHttpClient;
import okhttp3.FormBody;
import okhttp3.Request;
import okhttp3.MultipartBody;
import okhttp3.RequestBody;
import okhttp3.ResponseBody;
import okhttp3.Response;
import okhttp3.MediaType;

import okio.Buffer;
import okio.BufferedSink;
import okio.BufferedSource;
import okio.Okio;

import java.io.File;
import java.io.FileOutputStream;

class Steps {
    Step[] steps;
}

class Step {
    long id;
    String[] image_urls;
}

class ImageRequest {
    String url;
    int rotate;
    int flip;
    String crop;
}


public class ClientPhotolab {
    public static final String API_ENDPOINT = "http://api-soft.photolab.me";
    public static final String API_UPLOAD_ENDPOINT = "http://upload-soft.photolab.me/upload.php";
    public static final String API_ENDPOINT_PROXY = "http://api-proxy-soft.photolab.me";


    public Steps PhotolabStepsAdvanced(long comboId) throws Exception {
        RequestBody formBody = new FormBody.Builder()
                .add("combo_id", String.format("%d", comboId))
                .build();
        String endpoint = String.format("%s/photolab_steps_advanced.php", API_ENDPOINT);

        String data = this.query(endpoint, formBody);
        Gson gson = new Gson();
        Steps st = gson.fromJson(data, Steps.class);
        return st;
    }

    public String TemplateProcess(String template_name, ImageRequest[] ims) throws Exception {
        FormBody.Builder formBuilder = new FormBody.Builder();
        formBuilder.add("template_name", template_name);
        for (int i = 0; i < ims.length; i++) {
            ImageRequest im = ims[i];
            String t = String.format("%d", i + 1);
            formBuilder.add(String.format("image_url[%s]", t), im.url);
            formBuilder.add(String.format("rotate[%s]", t), Integer.toString(im.rotate));
            formBuilder.add(String.format("crop[%s]", t), im.crop);
            formBuilder.add(String.format("flip[%s]", t), Integer.toString(im.flip));
        }
        String endpoint = String.format("%s/template_process.php", API_ENDPOINT);
        Gson gson = new Gson();
        return this.query(endpoint, formBuilder.build());
    }

    public String PhotolabProcess(String template_name, ImageRequest[] ims) throws Exception {
        return TemplateProcess(template_name, ims);
    }

    public void DownloadFile(String url, String dest) throws Exception {
        OkHttpClient client = new OkHttpClient();
        Request request = new Request.Builder().url(url).build();
        Response response = client.newCall(request).execute();
        ResponseBody body = response.body();
        BufferedSource source = body.source();
        FileOutputStream destFile = new FileOutputStream(dest);
        BufferedSink sink = Okio.buffer(Okio.sink(destFile));
        Buffer sinkBuffer = sink.buffer();
        int bufferSize = 64 * 1024;
        for (; (source.read(sinkBuffer, bufferSize)) != -1; ) {
            sink.emit();
        }
        sink.flush();
        sink.close();
        source.close();
        return;
    }

    protected String query(String url, RequestBody formBody) throws Exception {
        OkHttpClient client = new OkHttpClient();
        Request request = new Request.Builder()
                .url(url)
                .post(formBody)
                .build();
        Response response = client.newCall(request).execute();
        if (!response.isSuccessful()) {
            String message = String.format("StatusCode %d from: %s", response.code(), request.url()) + " Response: " + response.body().string();
            throw new Exception(message);
        }
        return response.body().string();
    }

    public String ImageUpload(String image) throws Exception {
        String url = API_UPLOAD_ENDPOINT;
        File file = new File(image);
        RequestBody formBody = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("file1", file.getName(),
                        RequestBody.create(MediaType.parse("multipart/form-data"), file))
                .addFormDataPart("no_resize", "1")
                .build();

        return query(url, formBody);
    }

    public String TemplateUpload(String resources) throws Exception {
        String url = String.format("%s/template_upload.php", API_ENDPOINT_PROXY);
        File file = new File(resources);
        RequestBody formBody = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("resources", file.getName(),
                        RequestBody.create(MediaType.parse("multipart/form-data"), file))
                .build();
        return query(url, formBody);
    }
}
