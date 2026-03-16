package demo;

import org.deeplearning4j.nn.conf.MultiLayerConfiguration;
import org.deeplearning4j.nn.conf.NeuralNetConfiguration;
import org.deeplearning4j.nn.conf.layers.OutputLayer;
import org.deeplearning4j.nn.multilayer.MultiLayerNetwork;
import org.nd4j.linalg.activations.Activation;
import org.nd4j.linalg.api.ndarray.INDArray;
import org.nd4j.linalg.dataset.DataSet;
import org.nd4j.linalg.factory.Nd4j;
import org.nd4j.linalg.learning.config.Sgd;
import org.nd4j.linalg.lossfunctions.LossFunctions;

public class LinearDemo {
    public static void main(String[] args) {
        System.out.println("Backend: " + Nd4j.getBackend().getClass().getSimpleName());

        // y = 2x + 1
        INDArray X = Nd4j.create(new float[][]{{1}, {2}, {3}, {4}});
        INDArray y = Nd4j.create(new float[][]{{3}, {5}, {7}, {9}});
        DataSet data = new DataSet(X, y);

        MultiLayerConfiguration conf = new NeuralNetConfiguration.Builder()
                .updater(new Sgd(0.01))
                .list()
                .layer(new OutputLayer.Builder(LossFunctions.LossFunction.MSE)
                        .nIn(1).nOut(1)
                        .activation(Activation.IDENTITY)
                        .build())
                .build();

        MultiLayerNetwork model = new MultiLayerNetwork(conf);
        model.init();

        for (int epoch = 0; epoch < 500; epoch++) {
            model.fit(data);
        }

        double w = model.getLayer(0).getParam("W").getDouble(0);
        double b = model.getLayer(0).getParam("b").getDouble(0);
        System.out.printf("Learned: y = %.2fx + %.2f  (loss: %.6f)%n",
                w, b, model.score());
    }
}
