import { ResponseContext, RequestContext, HttpFile, HttpInfo } from '../http/http';
import { Configuration, PromiseConfigurationOptions, wrapOptions } from '../configuration'
import { PromiseMiddleware, Middleware, PromiseMiddlewareWrapper } from '../middleware';

import { Aes128 } from '../models/Aes128';
import { Aes256 } from '../models/Aes256';
import { Channel } from '../models/Channel';
import { ChannelConfiguration } from '../models/ChannelConfiguration';
import { ChannelProfile } from '../models/ChannelProfile';
import { ChannelState } from '../models/ChannelState';
import { ChannelStatus } from '../models/ChannelStatus';
import { ChannelType } from '../models/ChannelType';
import { ConfigurationData } from '../models/ConfigurationData';
import { ConnectRequestContent } from '../models/ConnectRequestContent';
import { ConnectResponseContent } from '../models/ConnectResponseContent';
import { Connection } from '../models/Connection';
import { Critical } from '../models/Critical';
import { Degraded } from '../models/Degraded';
import { DeprovisionResponseContent } from '../models/DeprovisionResponseContent';
import { DeviceConfiguration } from '../models/DeviceConfiguration';
import { DeviceRegistration } from '../models/DeviceRegistration';
import { DeviceStatus } from '../models/DeviceStatus';
import { DisconnectResponseContent } from '../models/DisconnectResponseContent';
import { EncryptionAes } from '../models/EncryptionAes';
import { EncryptionAes128 } from '../models/EncryptionAes128';
import { EncryptionAes256 } from '../models/EncryptionAes256';
import { EnumValues } from '../models/EnumValues';
import { ErrorDetails } from '../models/ErrorDetails';
import { GetConfigurationResponseContent } from '../models/GetConfigurationResponseContent';
import { GetConnectionStatusResponseContent } from '../models/GetConnectionStatusResponseContent';
import { Health } from '../models/Health';
import { Healthy } from '../models/Healthy';
import { IdAndValue } from '../models/IdAndValue';
import { Profile } from '../models/Profile';
import { ProfileDefinition } from '../models/ProfileDefinition';
import { RangeValues } from '../models/RangeValues';
import { ReportActualConfigurationRequestContent } from '../models/ReportActualConfigurationRequestContent';
import { ReportActualConfigurationResponseContent } from '../models/ReportActualConfigurationResponseContent';
import { ReportStatusRequestContent } from '../models/ReportStatusRequestContent';
import { ReportStatusResponseContent } from '../models/ReportStatusResponseContent';
import { RistCaller } from '../models/RistCaller';
import { RistCallerTransportProtocol } from '../models/RistCallerTransportProtocol';
import { RistListener } from '../models/RistListener';
import { RistListenerTransportProtocol } from '../models/RistListenerTransportProtocol';
import { RistStreamIdentifier } from '../models/RistStreamIdentifier';
import { Rtp } from '../models/Rtp';
import { RtpFecConfiguration } from '../models/RtpFecConfiguration';
import { RtpFecStreamConfig } from '../models/RtpFecStreamConfig';
import { RtpTransportProtocol } from '../models/RtpTransportProtocol';
import { Setting } from '../models/Setting';
import { SettingsChoice } from '../models/SettingsChoice';
import { SrtCaller } from '../models/SrtCaller';
import { SrtCallerTransportProtocol } from '../models/SrtCallerTransportProtocol';
import { SrtListener } from '../models/SrtListener';
import { SrtListenerTransportProtocol } from '../models/SrtListenerTransportProtocol';
import { StandardSettings } from '../models/StandardSettings';
import { StatusValue } from '../models/StatusValue';
import { StreamId } from '../models/StreamId';
import { SynchronizationSource } from '../models/SynchronizationSource';
import { Thumbnail } from '../models/Thumbnail';
import { TransportProtocol } from '../models/TransportProtocol';
import { TransportProtocolName } from '../models/TransportProtocolName';
import { UnhealthyStateDescription } from '../models/UnhealthyStateDescription';
import { ZixiPull } from '../models/ZixiPull';
import { ZixiPullTransportProtocol } from '../models/ZixiPullTransportProtocol';
import { ZixiPush } from '../models/ZixiPush';
import { ZixiPushTransportProtocol } from '../models/ZixiPushTransportProtocol';
import { ObservableDefaultApi } from './ObservableAPI';

import { DefaultApiRequestFactory, DefaultApiResponseProcessor} from "../apis/DefaultApi";
export class PromiseDefaultApi {
    private api: ObservableDefaultApi

    public constructor(
        configuration: Configuration,
        requestFactory?: DefaultApiRequestFactory,
        responseProcessor?: DefaultApiResponseProcessor
    ) {
        this.api = new ObservableDefaultApi(configuration, requestFactory, responseProcessor);
    }

    /**
     * @param connectRequestContent
     */
    public connectWithHttpInfo(connectRequestContent: ConnectRequestContent, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ConnectResponseContent>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.connectWithHttpInfo(connectRequestContent, observableOptions);
        return result.toPromise();
    }

    /**
     * @param connectRequestContent
     */
    public connect(connectRequestContent: ConnectRequestContent, _options?: PromiseConfigurationOptions): Promise<ConnectResponseContent> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.connect(connectRequestContent, observableOptions);
        return result.toPromise();
    }

    /**
     * @param hostId
     * @param [force]
     */
    public deprovisionWithHttpInfo(hostId: string, force?: boolean, _options?: PromiseConfigurationOptions): Promise<HttpInfo<DeprovisionResponseContent>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.deprovisionWithHttpInfo(hostId, force, observableOptions);
        return result.toPromise();
    }

    /**
     * @param hostId
     * @param [force]
     */
    public deprovision(hostId: string, force?: boolean, _options?: PromiseConfigurationOptions): Promise<DeprovisionResponseContent> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.deprovision(hostId, force, observableOptions);
        return result.toPromise();
    }

    /**
     */
    public disconnectWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<DisconnectResponseContent>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.disconnectWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     */
    public disconnect(_options?: PromiseConfigurationOptions): Promise<DisconnectResponseContent> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.disconnect(observableOptions);
        return result.toPromise();
    }

    /**
     */
    public getConfigurationWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<GetConfigurationResponseContent>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.getConfigurationWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     */
    public getConfiguration(_options?: PromiseConfigurationOptions): Promise<GetConfigurationResponseContent> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.getConfiguration(observableOptions);
        return result.toPromise();
    }

    /**
     */
    public getConnectionStatusWithHttpInfo(_options?: PromiseConfigurationOptions): Promise<HttpInfo<GetConnectionStatusResponseContent>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.getConnectionStatusWithHttpInfo(observableOptions);
        return result.toPromise();
    }

    /**
     */
    public getConnectionStatus(_options?: PromiseConfigurationOptions): Promise<GetConnectionStatusResponseContent> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.getConnectionStatus(observableOptions);
        return result.toPromise();
    }

    /**
     * @param reportActualConfigurationRequestContent
     */
    public reportActualConfigurationWithHttpInfo(reportActualConfigurationRequestContent: ReportActualConfigurationRequestContent, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ReportActualConfigurationResponseContent>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.reportActualConfigurationWithHttpInfo(reportActualConfigurationRequestContent, observableOptions);
        return result.toPromise();
    }

    /**
     * @param reportActualConfigurationRequestContent
     */
    public reportActualConfiguration(reportActualConfigurationRequestContent: ReportActualConfigurationRequestContent, _options?: PromiseConfigurationOptions): Promise<ReportActualConfigurationResponseContent> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.reportActualConfiguration(reportActualConfigurationRequestContent, observableOptions);
        return result.toPromise();
    }

    /**
     * @param reportStatusRequestContent
     */
    public reportStatusWithHttpInfo(reportStatusRequestContent: ReportStatusRequestContent, _options?: PromiseConfigurationOptions): Promise<HttpInfo<ReportStatusResponseContent>> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.reportStatusWithHttpInfo(reportStatusRequestContent, observableOptions);
        return result.toPromise();
    }

    /**
     * @param reportStatusRequestContent
     */
    public reportStatus(reportStatusRequestContent: ReportStatusRequestContent, _options?: PromiseConfigurationOptions): Promise<ReportStatusResponseContent> {
        const observableOptions = wrapOptions(_options);
        const result = this.api.reportStatus(reportStatusRequestContent, observableOptions);
        return result.toPromise();
    }


}



